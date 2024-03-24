tailwind.config = {
    theme: {
        fontFamily: {
            sans: ['Ubuntu'],
        }
    }
};

Chart.defaults.font.family = "Ubuntu";
Chart.defaults.color = "rgb(55, 65, 81)";

const colors = {
    blue: "#3b82f6",
    yellow: "#eab308",
    red: "#ef4444",
    green: "#22c55e",
    orange: "#f97316",
}

function parseData(id, key) {
    return JSON.parse(document.getElementById(id).dataset[key]);
}

function makeDataset(label, color, data) {
    return {
        label: label,
        data: data,
        tension: 0.3,
        borderColor: color,
        pointRadius: 7,
    };
}

function newLineChart(datasets) {
    return {
        type: "line",
        data: {
            datasets: datasets,
        },
        options: {
            plugins: {
                legend: {
                    position: 'bottom',
                },
            },
        },
    };
}

function newBarChart(datasets) {
    return {
        type: "bar",
        data: {
            datasets: datasets,
        },
        options: {
            indexAxis: 'y',
            scales: {
                y: {
                    ticks: {
                        autoSkip: false,
                    },
                },
            },
            plugins: {
                legend: {
                    position: 'bottom',
                },
            },
        },
    };
}

new Chart(
    document.getElementById("users-activity"),
    newLineChart([
        makeDataset(
            "количество сообщений",
            colors.blue,
            parseData("users-activity", 'messageCount'),
        ),
        makeDataset(
            "количество команд",
            colors.yellow,
            parseData("users-activity", 'commandCount'),
        ),
    ]),
);

new Chart(
    document.getElementById("top-commands"),
    newBarChart([
        {
            label: "количество вызовов",
            data: parseData("top-commands", 'topCommands'),
            backgroundColor: colors.red,
        },
    ]),
);

new Chart(
    document.getElementById("top-users"),
    newBarChart([
        {
            label: "количество сообщений",
            data: parseData("top-users", 'topUsers'),
            backgroundColor: colors.green,
        },
    ]),
);
