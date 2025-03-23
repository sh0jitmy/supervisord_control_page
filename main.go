package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, htmlContent)
	})

	fmt.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

const htmlContent = `<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <title>Supervisord Process Control</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100 p-8">
    <div class="max-w-4xl mx-auto">
        <h1 class="text-3xl font-bold mb-8">Supervisord プロセス管理</h1>

        <div id="process-list" class="space-y-4"></div>
    </div>

    <script>
        if (!window.apiUrl) {
            window.apiUrl = '/RPC2';
        }

        async function fetchProcesses() {
            const response = await fetch(window.apiUrl, {
                method: 'POST',
                headers: { 'Content-Type': 'text/xml' },
                body: `<?xml version="1.0"?><methodCall><methodName>supervisor.getAllProcessInfo</methodName></methodCall>`
            });
            const xml = await response.text();
            const parser = new DOMParser();
            const xmlDoc = parser.parseFromString(xml, 'text/xml');

            const processes = Array.from(xmlDoc.getElementsByTagName('struct')).map(parseProcess);

            renderProcesses(processes);
        }

        function parseProcess(struct) {
            const getValue = (name) => struct.querySelector(`name:contains('${name}')`).nextElementSibling.textContent;
            return {
                name: getValue('name'),
                state: getValue('statename')
            };
        }

        async function controlProcess(action, name) {
            const method = action === 'start' ? 'supervisor.startProcess' : 'supervisor.stopProcess';
            await fetch(window.apiUrl, {
                method: 'POST',
                headers: { 'Content-Type': 'text/xml' },
                body: `<?xml version="1.0"?><methodCall><methodName>${method}</methodName><params><param><value><string>${name}</string></value></param></params></methodCall>`
            });
            fetchProcesses();
        }

        function renderProcesses(processes) {
            const container = document.getElementById('process-list');
            container.innerHTML = '';

            processes.forEach(proc => {
                const isRunning = proc.state === 'RUNNING';

                const div = document.createElement('div');
                div.className = 'p-4 bg-white rounded-lg shadow';

                div.innerHTML = `
                    <h2 class="text-xl font-semibold">${proc.name}</h2>
                    <p class="text-gray-700">状態: <span class="font-mono">${proc.state}</span></p>
                    <button onclick="controlProcess('${isRunning ? 'stop' : 'start'}', '${proc.name}')"
                        class="mt-4 px-4 py-2 text-white rounded ${isRunning ? 'bg-red-500 hover:bg-red-600' : 'bg-green-500 hover:bg-green-600'}">
                        ${isRunning ? '停止' : '起動'}
                    </button>
                `;

                container.appendChild(div);
            });
        }

        fetchProcesses();
    </script>
</body>
</html>`
