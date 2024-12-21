async function fetchData() {
    try {
        const response = await fetch('http://localhost:8092/v1/dnt/table');
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('There has been a problem with your fetch operation:', error);
    }
}