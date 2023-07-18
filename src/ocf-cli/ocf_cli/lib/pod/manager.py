import os
import shutil
import signal
import subprocess
import tempfile

try:
    from importlib.metadata import version
except ModuleNotFoundError:
    # Fallback for Python 3.7
    from importlib_metadata import version

from itertools import filterfalse
from pathlib import Path
from typing import Dict, List, Optional

import psutil
from rich import box
from rich.console import Console
from rich.panel import Panel
from rich.table import Table
from rich.text import Text

from ocf_cli.lib.pod import config
from ocf_cli.lib.pod import Pod
from ocf_cli.lib.pod.utils import kill_proc_tree, logger, wait_created

console = Console(highlight=False)

class PodManager(object):
    def __init__(self, tcmpod_dir: Path = None):
        default_dir = Path(tempfile.gettempdir()) / "ocf_cli"
        self._tcmpod_dir = tcmpod_dir or default_dir
        logger.debug(f"Initialized within {self._tcmpod_dir} dir")
        if not self._tcmpod_dir.exists():
            self._tcmpod_dir.mkdir(parents=True, exist_ok=True)

    @staticmethod
    def stats(pods: List[Pod], verbose: bool = False):
        if not pods:
            console.print(
                f"{config.ICON_INFO} No pods are currently running",
                style=f"{config.COLOR_MAIN} bold",
            )
            return

        package_version = version(__package__)
        table = Table(
            show_header=True,
            header_style=f"{config.COLOR_MAIN} bold",
            box=box.HEAVY_EDGE,
            caption_style="dim",
            caption_justify="right",
        )
        table.add_column("#", style="dim", width=2)
        table.add_column("Name")
        table.add_column("PID")
        if verbose:
            table.add_column("Command")
        table.add_column("Status")
        table.add_column("RC", justify="right")
        table.add_column("Runtime", justify="right")

        active_pods = 0
        for pod in pods:
            active_pods += 1 if pod.active else 0
            pid_text = f"{pod.pid}" if pod.active else Text(f"{pod.pid}", style="dim")
            command_text = Text(
                pod.cmd, overflow="ellipsis", style=f"{config.COLOR_ACCENT}"
            )
            status_text = PodManager._get_status_text(pod.status)
            command_text.truncate(config.TRUNCATE_LENGTH)
            row = [
                f"{pod.hid}",
                pod.name,
                pid_text,
                command_text if verbose else None,
                status_text,
                f"{pod.rc}" if pod.rc is not None else "",
                pod.runtime,
            ]
            table.add_row(*filterfalse(lambda x: x is None, row))

        if verbose:
            table.title = f"{config.ICON_POD} {__package__}, {package_version}"
            table.caption = f"{active_pods} active / {len(pods)} total"
        console.print(table)

    @staticmethod
    def _get_status_text(status) -> Text:
        color = config.STATUS_COLORS.get(status)
        status_text = Text()
        status_text.append(config.ICON_STATUS, style=color)
        status_text.append(f" {status}")
        return status_text

    @staticmethod
    def show(pod: Pod, verbose: bool = False):
        status_table = Table(show_header=False, show_footer=False, box=box.SIMPLE)

        status_text = PodManager._get_status_text(pod.status)
        status_table.add_row("Status:", status_text)

        status_table.add_row("PID:", f"{pod.pid}")

        if pod.rc is not None:
            status_table.add_row("Return code:", f"{pod.rc}")

        cmd_text = Text(f"{pod.cmd}", style=f"{config.COLOR_ACCENT} bold")
        status_table.add_row("Command:", cmd_text)

        proc = pod.proc
        if verbose and proc is not None:
            status_table.add_row("Working dir:", f"{proc.cwd()}")
            status_table.add_row("Parent PID:", f"{proc.ppid()}")
            status_table.add_row("User:", f"{proc.username()}")

        if verbose:
            status_table.add_row("Stdout file:", f"{pod.stdout_path}")
            status_table.add_row("Stderr file:", f"{pod.stderr_path}")

        start_time = pod.start_time
        end_time = pod.end_time
        if verbose and start_time:
            status_table.add_row("Start time:", f"{start_time}")

        if verbose and end_time:
            status_table.add_row("End time:", f"{end_time}")

        status_table.add_row("Runtime:", f"{pod.runtime}")

        status_panel = Panel(
            status_table,
            expand=verbose,
            title=f"Pod {config.ICON_POD}{pod.hid}",
            subtitle=pod.name,
        )
        console.print(status_panel)

        environ = pod.env
        if verbose and environ is not None:
            env_table = Table(show_header=False, show_footer=False, box=None)
            env_table.add_column("", justify="right")
            env_table.add_column("", justify="left", style=config.COLOR_ACCENT)

            for key, value in environ.items():
                env_table.add_row(key, Text(value, overflow="fold"))

            env_panel = Panel(
                env_table,
                title="Environment",
                subtitle=f"{len(environ)} items",
                border_style=config.COLOR_MAIN,
            )
            console.print(env_panel)

    @property
    def dir(self) -> Path:
        return self._tcmpod_dir

    def _get_pod_dirs(self) -> List[str]:
        pod_dirs = filter(str.isdigit, os.listdir(self._tcmpod_dir))
        return sorted(pod_dirs, key=int)

    def _get_pod_names_map(self) -> Dict[str, str]:
        names = {}
        for dir in self._get_pod_dirs():
            filename = self._tcmpod_dir / dir / "name"
            if filename.exists():
                with open(filename) as f:
                    name = f.read().strip()
                    names[name] = dir
        return names

    def get_next_pod_id(self):
        dirs = self._get_pod_dirs()
        return 1 if not dirs else int(dirs[-1]) + 1

    def get_pod(self, pod_alias: str) -> Optional[Pod]:
        dirs = self._get_pod_dirs()
        # Check by pod id
        if pod_alias in dirs:
            return Pod(self._tcmpod_dir / pod_alias)

        # Check by pod name
        names_map = self._get_pod_names_map()
        if pod_alias in names_map:
            return Pod(self._tcmpod_dir / names_map[pod_alias])

    def get_pods(self) -> List[Pod]:
        pods = []
        if not self._tcmpod_dir.exists():
            return pods

        for dir in self._get_pod_dirs():
            pod_path = self._tcmpod_dir / dir
            pods.append(Pod(pod_path=pod_path))
        return pods

    def create_pod(
        self,
        cmd: str,
        name: Optional[str] = None,
        additional_env: Optional[dict] = None,
    ) -> Pod:
        hid = self.get_next_pod_id()
        pod_dir = self._tcmpod_dir / f"{hid}"
        pod_dir.mkdir()
        return Pod(pod_dir, cmd=cmd, name=name, additional_env=additional_env)

    def run_pod(self, pod: Pod):
        with open(pod.stdout_path, "w") as stdout_pipe, open(
            pod.stderr_path, "w"
        ) as stderr_pipe:
            console.print(f"{config.ICON_INFO} Launching", pod)
            proc = subprocess.Popen(
                pod.cmd,
                shell=True,
                stdout=stdout_pipe,
                stderr=stderr_pipe,
                env=pod._env,
            )
            pid = proc.pid
            logger.debug(f"Attaching pod {pod} to pid {pid}")

            pod.attach(pid)

    def _check_fast_failure(self, pod: Pod):
        if wait_created(pod._rc_file) and pod.rc != 0:
            console.print(
                f"{config.ICON_INFO} Pod exited too quickly. stderr message:",
                style=f"{config.COLOR_ERROR} bold",
            )
            with open(pod.stderr_path) as f:
                console.print(f.read())

    def pause_pod(self, pod: Pod):
        proc = pod.proc
        if proc is not None:
            proc.suspend()
            console.print(f"{config.ICON_INFO} Paused", pod)
        else:
            console.print(
                f"{config.ICON_INFO} Cannot pause. Pod {pod} is not running",
                style=f"{config.COLOR_ERROR} bold",
            )

    def resume_pod(self, pod: Pod):
        proc = pod.proc
        if proc is not None and proc.status() == psutil.STATUS_STOPPED:
            proc.resume()
            console.print(f"{config.ICON_INFO} Resumed", pod)
        else:
            console.print(
                f"{config.ICON_INFO} Cannot resume. Hap {pod} is not suspended",
                style=f"{config.COLOR_ERROR} bold",
            )

    def run(
        self,
        cmd: str,
        name: Optional[str] = None,
        check: bool = False,
        additional_env: Optional[dict] = None,
    ):
        pod = self.create_pod(cmd=cmd, name=name, additional_env=additional_env)
        pid = os.fork()
        if pid == 0:
            self.run_pod(pod)
        else:
            if check:
                self._check_fast_failure(pod)

    def logs(self, pod: Pod, stderr: bool = False, follow: bool = False):
        filepath = pod.stderr_path if stderr else pod.stdout_path
        if follow:
            console.print(
                f"{config.ICON_INFO} Streaming {filepath} file...",
                style=f"{config.COLOR_MAIN} bold",
            )
            return subprocess.run(["tail", "-f", filepath])
        else:
            return subprocess.run(["cat", filepath])

    def clean(self, skip_failed: bool = False):
        def to_clean(pod):
            if pod.rc is not None:
                return pod.rc == 0 or not skip_failed
            return False

        pods = list(filter(to_clean, self.get_pods()))
        for pod in pods:
            logger.debug(f"Removing {pod.path}")
            shutil.rmtree(pod.path)

        if pods:
            console.print(
                f"{config.ICON_INFO} Deleted {len(pods)} finished pods",
                style=f"{config.COLOR_MAIN} bold",
            )
        else:
            console.print(
                f"{config.ICON_INFO} Nothing to clean",
                style=f"{config.COLOR_ERROR} bold",
            )

    def clean_pod(self, pod: Pod):
        shutil.rmtree(pod.path)

    def kill(self, pods: List[Pod]):
        killed_counter = 0
        for pod in pods:
            if pod.active:
                logger.info(f"Killing {pod}...")
                kill_proc_tree(pod.pid)
                killed_counter += 1

        if killed_counter:
            console.print(
                f"{config.ICON_KILLED} Killed {killed_counter} active pods",
                style=f"{config.COLOR_MAIN} bold",
            )
        else:
            console.print(
                f"{config.ICON_INFO} No active pods to kill",
                style=f"{config.COLOR_ERROR} bold",
            )

    def signal(self, pod: Pod, sig: signal.Signals):
        if pod.active:
            sig_text = (
                f"[bold]{sig.name}[/] ([{config.COLOR_MAIN}]{signal.strsignal(sig)}[/])"
            )
            console.print(f"{config.ICON_INFO} Sending {sig_text} to pod {pod}")
            pod.proc.send_signal(sig)
        else:
            console.print(
                f"{config.ICON_INFO} Cannot send signal to the inactive pod",
                style=f"{config.COLOR_ERROR} bold",
            )
