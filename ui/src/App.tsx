import { MouseEvent, useRef, useState } from 'react'
import { Button } from 'primereact/button'
import { Message } from 'primereact/message'
import CodeMirror from '@uiw/react-codemirror';
import { html } from '@codemirror/lang-html';
import { javascript } from '@codemirror/lang-javascript';
import { json } from '@codemirror/lang-json';
import { vscodeDark } from '@uiw/codemirror-theme-vscode'
import './App.css'

import { ExampleHTML, ExampleData } from './examples'

function App() {
  const [dataCode, setDataCode] = useState(ExampleData)
  const [templateCode, setTemplateCode] = useState(ExampleHTML)
  const [editorTheme, setEditorTheme] = useState(vscodeDark)

  const downloadRef: any = useRef(null)

  const [downloadName, setDownloadName] = useState('test.pdf')


  async function uploadAndRenderPDF(event: MouseEvent<HTMLButtonElement>): Promise<boolean> {
    event.preventDefault()

    const templateRequest = {
      Name: 'test.html',
      Content: templateCode
    }

    let response = await fetch("/_studio/upload-template", {
      method: 'POST',
      mode: 'cors',
      cache: 'no-cache',
      credentials: 'same-origin',
      headers: {
        'Content-Type': 'application/json',
        'X-Greypot-Studio-Version': '0.0.1-dev',
      },
      redirect: 'error',
      referrerPolicy: 'no-referrer',
      body: JSON.stringify(templateRequest),
    });

    if (response.ok) {
      let testDataJSON = JSON.parse(dataCode)
      let response = await fetch(`/_studio/reports/export/pdf/${templateRequest.Name}`, {
        method: 'POST',
        mode: 'cors',
        cache: 'no-cache',
        credentials: 'same-origin',
        headers: {
          'Content-Type': 'application/json',
          'X-Greypot-Studio-Version': '0.0.1-dev',
        },
        redirect: 'error',
        referrerPolicy: 'no-referrer',
        body: JSON.stringify(testDataJSON),
      });

      if (response.ok) {
        type ExportResponse = {
          data: string,
          type: string,
          reportId: string
        }
        let res = await response.json() as ExportResponse;

        downloadRef.current.setAttribute("href", `data:application/octet-stream;base64,${res.data}`)
        downloadRef.current.setAttribute("download", templateRequest.Name.replace(".html", ".pdf"))
        await downloadRef.current.click();
      }
    }

    return false
  }

  return (
    <div className="App">
      <div className="masthead">
        <div className="grid">
          <div className="col-10">
            <h1>Greypot Studio v0.0.1</h1>
          </div>
          <div className="col-2">
            <Message severity="warn" text="Still in development" />
          </div>
        </div>
      </div>
      <div className="grid grid-nogutter">
        <div className="col-8">
          <h2>HTML Design Template</h2>
          <CodeMirror
            width="100%"
            height="400px"
            extensions={[html(), javascript()]}
            value={templateCode}
            onChange={(e) => setTemplateCode(e)}
            theme={editorTheme}
          // options={options}
          // editorDidMount={editorDidMount}
          />
        </div>
        <div className="col-4">
          <h2>Test Data</h2>
          {/* <MonacoEditor
            width="100%"
            height="400"
            language="json"
            theme={editorTheme}
            value={dataCode}
            options={options}
            onChange={(newValue, e) => setDataCode(newValue)}
          /> */}
          <CodeMirror
            width="100%"
            height="400px"
            extensions={[json()]}
            theme={editorTheme}
            value={dataCode}
            onChange={(e) => setDataCode(e)}
          />
        </div>
      </div>

      <div className="action-area p-3">
        <a style={{ display: 'none' }} ref={downloadRef} download={downloadName}></a>
        <Button label="PDF Preview with Test Data" onClick={uploadAndRenderPDF} />
      </div>

      <footer>
        <p>👋 Mulibwanji</p>
        <a href="https://nndi.cloud/oss/greypot">Greypot Studio</a> is an open-source project brought to you from Malawi by NNDI.
      </footer>
    </div>
  )
}

export default App
