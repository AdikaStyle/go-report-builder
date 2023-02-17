import { MouseEvent, useRef, useState } from 'react'
import { Button } from 'primereact/button'
import { Message } from 'primereact/message'
import MonacoEditor from 'react-monaco-editor';
import { editor } from 'monaco-editor/esm/vs/editor/editor.api';
import './App.css'
import {ExampleHTML, ExampleData} from './examples'

function App() {
  const [dataCode, setDataCode] = useState(ExampleData)
  const [templateCode, setTemplateCode] = useState(ExampleHTML)

  const downloadRef: any = useRef(null)

  const [downloadName, setDownloadName] = useState('test.pdf')

  let editorDidMount = (editor: any, monaco: any) => {
    editor.focus()
  }

  let onChange = (newValue: string, e: editor.IModelContentChangedEvent) => {
    setTemplateCode(newValue);
  }
  const options = {
    selectOnLineNumbers: true
  };

  async function uploadAndRenderPDF(event: MouseEvent<HTMLButtonElement>) : Promise<boolean> {
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
          <MonacoEditor
            width="100%"
            height="400"
            language="html"
            theme="vs-dark"
            value={templateCode}
            options={options}
            onChange={onChange}
            editorDidMount={editorDidMount}
          />
        </div>
        <div className="col-4">
          <h2>Test Data</h2>
          <MonacoEditor
            width="100%"
            height="400"
            language="json"
            theme="vs-dark"
            value={dataCode}
            options={options}
            onChange={(newValue, e) => setDataCode(newValue)}
          />
        </div>
      </div>

      <div className="action-area p-3">
        <a style={{ display: 'none'}} ref={downloadRef} download={downloadName}></a>
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
