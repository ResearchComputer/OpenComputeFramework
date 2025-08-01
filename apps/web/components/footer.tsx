import { Separator } from "@/components/ui/separator";
import {
  DribbbleIcon,
  GithubIcon,
  TwitchIcon,
  TwitterIcon,
} from "lucide-react";
import Link from "next/link";

const footerSections = [
  {
    title: "Framework",
    links: [
      {
        title: "Documentation",
        href: "https://github.com/ResearchComputer/OpenComputeFramework",
      },
      {
        title: "Getting Started",
        href: "#",
      },
      {
        title: "API Reference",
        href: "https://api.research.computer",
      },
      {
        title: "Examples",
        href: "#",
      },
      {
        title: "Releases",
        href: "https://github.com/ResearchComputer/OpenComputeFramework/releases",
      },
      {
        title: "Docker Images",
        href: "https://ghcr.io/xiaozheyao/ocf",
      },
    ],
  },
  {
    title: "Community",
    links: [
      {
        title: "GitHub",
        href: "https://github.com/ResearchComputer/OpenComputeFramework",
      },
      {
        title: "Discussions",
        href: "https://github.com/ResearchComputer/OpenComputeFramework/discussions",
      },
      {
        title: "Issues",
        href: "https://github.com/ResearchComputer/OpenComputeFramework/issues",
      },
      {
        title: "Contributing",
        href: "https://github.com/ResearchComputer/OpenComputeFramework/blob/main/CONTRIBUTING.md",
      },
      {
        title: "Research Computer",
        href: "https://research.computer",
      },
    ],
  },
  {
    title: "Resources",
    links: [
      {
        title: "ML Inference API",
        href: "https://api.research.computer",
      },
      {
        title: "Triteia Backend",
        href: "https://api.research.computer/triteia/",
      },
      {
        title: "Inferencia Backend",
        href: "https://api.research.computer/inferencia/",
      },
      {
        title: "Network Status",
        href: "#",
      },
      {
        title: "LibP2P Docs",
        href: "https://libp2p.io",
      },
      {
        title: "CRDT Guide",
        href: "#",
      },
    ],
  },
  {
    title: "Deploy",
    links: [
      {
        title: "Docker",
        href: "https://hub.docker.com/r/xiaozheyao/ocf",
      },
      {
        title: "Binary Releases",
        href: "https://github.com/ResearchComputer/OpenComputeFramework/releases",
      },
      {
        title: "Build from Source",
        href: "https://github.com/ResearchComputer/OpenComputeFramework#building",
      },
      {
        title: "Kubernetes",
        href: "#",
      },
      {
        title: "Cloud Providers",
        href: "#",
      },
    ],
  },
  {
    title: "Legal",
    links: [
      {
        title: "Terms",
        href: "#",
      },
      {
        title: "Privacy",
        href: "#",
      },
      {
        title: "Cookies",
        href: "#",
      },
      {
        title: "Licenses",
        href: "#",
      },
      {
        title: "Settings",
        href: "#",
      },
      {
        title: "Contact",
        href: "#",
      },
    ],
  },
];

const Footer = () => {
  return (
    <footer className="mt-12 xs:mt-20 dark bg-background border-t">
      <div className="max-w-screen-xl mx-auto py-12 grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-7 gap-x-8 gap-y-10 px-6">
        <div className="col-span-full xl:col-span-2">
          {/* Logo */}
          <svg
            id="logo-7"
            width="124"
            height="32"
            viewBox="0 0 124 32"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M36.87 10.07H39.87V22.2H36.87V10.07ZM41.06 17.62C41.06 14.62 42.9 12.83 45.74 12.83C48.58 12.83 50.42 14.62 50.42 17.62C50.42 20.62 48.62 22.42 45.74 22.42C42.86 22.42 41.06 20.67 41.06 17.62ZM47.41 17.62C47.41 15.97 46.76 15 45.74 15C44.72 15 44.08 16 44.08 17.62C44.08 19.24 44.71 20.22 45.74 20.22C46.77 20.22 47.41 19.3 47.41 17.63V17.62ZM51.55 22.79H54.43C54.5671 23.0945 54.7988 23.3466 55.0907 23.5088C55.3826 23.6709 55.7191 23.7345 56.05 23.69C57.19 23.69 57.79 23.07 57.79 22.17V20.49H57.73C57.491 21.0049 57.1031 21.4363 56.6165 21.7287C56.1299 22.021 55.5668 22.1608 55 22.13C52.81 22.13 51.36 20.46 51.36 17.59C51.36 14.72 52.74 12.91 55.04 12.91C55.6246 12.8871 56.2022 13.0434 56.6955 13.3579C57.1888 13.6725 57.5742 14.1303 57.8 14.67V14.67V13H60.8V22.1C60.8 24.29 58.87 25.65 56.02 25.65C53.37 25.65 51.72 24.46 51.55 22.8V22.79ZM57.8 17.61C57.8 16.15 57.13 15.23 56.07 15.23C55.01 15.23 54.36 16.14 54.36 17.61C54.36 19.08 55 19.91 56.07 19.91C57.14 19.91 57.8 19.1 57.8 17.62V17.61ZM61.93 17.61C61.93 14.61 63.77 12.82 66.61 12.82C69.45 12.82 71.3 14.61 71.3 17.61C71.3 20.61 69.5 22.41 66.61 22.41C63.72 22.41 61.93 20.67 61.93 17.62V17.61ZM68.28 17.61C68.28 15.96 67.63 14.99 66.61 14.99C65.59 14.99 65 16 65 17.63C65 19.26 65.63 20.23 66.65 20.23C67.67 20.23 68.28 19.3 68.28 17.63V17.61ZM72.44 10.82C72.4321 10.5171 72.5144 10.2187 72.6763 9.96261C72.8383 9.70651 73.0726 9.50427 73.3496 9.38151C73.6266 9.25875 73.9338 9.221 74.2323 9.27305C74.5308 9.32511 74.8071 9.46462 75.0262 9.67389C75.2454 9.88317 75.3974 10.1528 75.4631 10.4486C75.5288 10.7444 75.5052 11.053 75.3952 11.3354C75.2853 11.6177 75.094 11.8611 74.8456 12.0346C74.5973 12.2081 74.3029 12.304 74 12.31C73.7992 12.3238 73.5977 12.2959 73.4082 12.2281C73.2186 12.1603 73.0452 12.0541 72.8987 11.916C72.7522 11.778 72.6358 11.6111 72.5569 11.4259C72.4779 11.2408 72.4381 11.0413 72.44 10.84V10.82ZM72.44 13.02H75.44V22.2H72.44V13.02ZM86.33 17.61C86.33 20.61 85 22.32 82.72 22.32C82.1354 22.3575 81.5533 22.2146 81.0525 21.9106C80.5517 21.6065 80.1564 21.156 79.92 20.62H79.86V25.14H76.86V13H79.86V14.64H79.92C80.1454 14.0951 80.5332 13.6329 81.0306 13.3162C81.528 12.9995 82.1109 12.8437 82.7 12.87C85 12.91 86.37 14.63 86.37 17.63L86.33 17.61ZM83.33 17.61C83.33 16.15 82.66 15.22 81.61 15.22C80.56 15.22 79.89 16.16 79.88 17.61C79.87 19.06 80.56 19.99 81.61 19.99C82.66 19.99 83.33 19.08 83.33 17.63V17.61ZM91.48 12.81C93.97 12.81 95.48 13.99 95.55 15.88H92.82C92.82 15.23 92.28 14.82 91.45 14.82C90.62 14.82 90.25 15.14 90.25 15.61C90.25 16.08 90.58 16.23 91.25 16.37L93.17 16.76C95 17.15 95.78 17.89 95.78 19.28C95.78 21.18 94.05 22.4 91.5 22.4C88.95 22.4 87.28 21.18 87.15 19.31H90.04C90.13 19.99 90.67 20.39 91.55 20.39C92.43 20.39 92.83 20.1 92.83 19.62C92.83 19.14 92.55 19.04 91.83 18.89L90.1 18.52C88.31 18.15 87.37 17.2 87.37 15.8C87.39 14 89 12.83 91.48 12.83V12.81ZM105.79 22.18H102.9V20.47H102.84C102.681 21.0441 102.331 21.5466 101.847 21.8941C101.363 22.2415 100.775 22.413 100.18 22.38C99.7242 22.4059 99.2682 22.3337 98.8427 22.1682C98.4172 22.0027 98.0322 21.7479 97.7137 21.4208C97.3952 21.0938 97.1505 20.7021 96.9964 20.2724C96.8422 19.8427 96.7821 19.3849 96.82 18.93V13H99.82V18.24C99.82 19.33 100.38 19.91 101.31 19.91C101.528 19.9104 101.744 19.8643 101.943 19.7746C102.141 19.6849 102.319 19.5537 102.463 19.3899C102.606 19.226 102.714 19.0333 102.777 18.8247C102.84 18.616 102.859 18.3962 102.83 18.18V13H105.83L105.79 22.18ZM107.24 13H110.14V14.77H110.2C110.359 14.2035 110.702 13.7057 111.174 13.3547C111.646 13.0037 112.222 12.8191 112.81 12.83C113.409 12.7821 114.003 12.9612 114.476 13.3318C114.948 13.7024 115.264 14.2372 115.36 14.83H115.42C115.601 14.2309 115.977 13.7093 116.488 13.3472C116.998 12.9851 117.615 12.8031 118.24 12.83C118.648 12.8163 119.054 12.8886 119.432 13.0422C119.811 13.1957 120.152 13.4272 120.435 13.7214C120.718 14.0157 120.936 14.3662 121.075 14.7501C121.213 15.134 121.27 15.5429 121.24 15.95V22.2H118.24V16.75C118.24 15.75 117.79 15.29 116.95 15.29C116.763 15.2884 116.577 15.327 116.406 15.4032C116.235 15.4794 116.082 15.5914 115.958 15.7317C115.834 15.872 115.741 16.0372 115.686 16.2163C115.631 16.3955 115.616 16.5843 115.64 16.77V22.2H112.79V16.71C112.79 15.79 112.34 15.29 111.52 15.29C111.331 15.2901 111.143 15.3303 110.971 15.408C110.798 15.4858 110.643 15.5993 110.518 15.741C110.392 15.8827 110.298 16.0495 110.241 16.2304C110.185 16.4112 110.167 16.6019 110.19 16.79V22.2H107.19L107.24 13Z"
              className="fill-foreground"
            />
            <path
              d="M28.48 10.62C27.9711 9.45636 27.2976 8.37193 26.48 7.4C25.2715 5.92034 23.7633 4.71339 22.0547 3.8586C20.3461 3.00382 18.4758 2.52057 16.567 2.44066C14.6582 2.36075 12.7541 2.68599 10.98 3.39499C9.20597 4.10398 7.60217 5.18065 6.2742 6.55413C4.94622 7.9276 3.92417 9.56675 3.27532 11.3637C2.62647 13.1606 2.36552 15.0746 2.50966 16.9796C2.65381 18.8847 3.19976 20.7376 4.1116 22.4164C5.02344 24.0953 6.28049 25.562 7.80001 26.72C8.77501 27.4779 9.85236 28.094 11 28.55C12.609 29.2094 14.3311 29.549 16.07 29.55C19.6594 29.5421 23.0992 28.1113 25.6355 25.5713C28.1717 23.0313 29.5974 19.5894 29.6 16C29.6026 14.1485 29.2213 12.3166 28.48 10.62V10.62ZM16.06 5.18999C17.6216 5.18983 19.1643 5.53113 20.58 6.18999V6.18999C20.2348 6.33916 19.8718 6.44335 19.5 6.5C18.2766 6.67709 17.1433 7.24507 16.2692 8.11917C15.3951 8.99326 14.8271 10.1266 14.65 11.35C14.5723 12.0361 14.2602 12.6744 13.7665 13.1572C13.2728 13.64 12.6277 13.9376 11.94 14C10.7166 14.1771 9.58327 14.7451 8.70918 15.6192C7.83509 16.4933 7.2671 17.6266 7.09001 18.85C7.03005 19.5024 6.7517 20.1155 6.30001 20.59V20.59C5.52066 18.9433 5.17056 17.1261 5.28228 15.3077C5.394 13.4893 5.96391 11.7287 6.93898 10.1897C7.91404 8.65079 9.26258 7.38351 10.8591 6.50584C12.4556 5.62817 14.2482 5.16864 16.07 5.16999L16.06 5.18999ZM7.79001 23C7.91001 22.89 8.03001 22.79 8.15001 22.67C9.03966 21.8075 9.61072 20.6689 9.77001 19.44C9.83459 18.7492 10.143 18.104 10.64 17.62C11.1183 17.1222 11.762 16.8163 12.45 16.76C13.6734 16.5829 14.8067 16.0149 15.6808 15.1408C16.5549 14.2667 17.1229 13.1334 17.3 11.91C17.3433 11.1875 17.6533 10.5068 18.17 10C18.6601 9.51185 19.3099 9.2171 20 9.16999C21.1239 9.01536 22.1721 8.51571 23 7.74C23.9427 8.52207 24.7413 9.46289 25.36 10.52C25.322 10.5713 25.2784 10.6183 25.23 10.66C24.7527 11.1622 24.1098 11.4748 23.42 11.54C22.1953 11.714 21.0603 12.281 20.1856 13.1556C19.311 14.0303 18.744 15.1653 18.57 16.39C18.4995 17.0784 18.1932 17.7213 17.703 18.2097C17.2127 18.6982 16.5687 19.0021 15.88 19.07C14.653 19.2457 13.5155 19.8126 12.6363 20.6863C11.7572 21.5601 11.1833 22.6941 11 23.92C10.9462 24.4087 10.7783 24.878 10.51 25.29C9.484 24.6808 8.5651 23.9072 7.79001 23V23ZM16.06 26.86C15.0453 26.8611 14.0354 26.7197 13.06 26.44C13.3937 25.818 13.6106 25.1401 13.7 24.44C13.7701 23.7531 14.075 23.1114 14.5632 22.6232C15.0514 22.135 15.6931 21.8301 16.38 21.76C17.6052 21.5849 18.7408 21.0178 19.6169 20.1435C20.4929 19.2693 21.0624 18.1348 21.24 16.91C21.3101 16.2231 21.615 15.5814 22.1032 15.0932C22.5914 14.605 23.2331 14.3001 23.92 14.23C24.842 14.1101 25.7208 13.7668 26.48 13.23C26.9016 14.8279 26.9515 16.5011 26.626 18.1213C26.3005 19.7415 25.6081 21.2657 24.6021 22.5768C23.5961 23.8878 22.3032 24.9511 20.8224 25.6849C19.3417 26.4187 17.7126 26.8036 16.06 26.81V26.86Z"
              className="fill-foreground"
            />
          </svg>

          <p className="mt-4 text-muted-foreground">
            Design amazing digital experiences that create more happy in the
            world.
          </p>
        </div>

        {footerSections.map(({ title, links }) => (
          <div key={title} className="xl:justify-self-end">
            <h6 className="font-semibold text-foreground">{title}</h6>
            <ul className="mt-6 space-y-4">
              {links.map(({ title, href }) => (
                <li key={title}>
                  <Link
                    href={href}
                    className="text-muted-foreground hover:text-foreground"
                  >
                    {title}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </div>
      <Separator />
      <div className="max-w-screen-xl mx-auto py-8 flex flex-col-reverse sm:flex-row items-center justify-between gap-x-2 gap-y-5 px-6">
        {/* Copyright */}
        <span className="text-muted-foreground text-center xs:text-start">
          &copy; {new Date().getFullYear()}{" "}
          <Link href="https://shadcnui-blocks.com" target="_blank">
            Shadcn UI Blocks
          </Link>
          . All rights reserved.
        </span>

        <div className="flex items-center gap-5 text-muted-foreground">
          <Link href="#" target="_blank">
            <TwitterIcon className="h-5 w-5" />
          </Link>
          <Link href="#" target="_blank">
            <DribbbleIcon className="h-5 w-5" />
          </Link>
          <Link href="#" target="_blank">
            <TwitchIcon className="h-5 w-5" />
          </Link>
          <Link href="#" target="_blank">
            <GithubIcon className="h-5 w-5" />
          </Link>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
