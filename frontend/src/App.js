import { useEffect, useState } from 'react';
import './App.css';

const apiBaseUrl = process.env.REACT_APP_API_URL || 'http://localhost:8080';

function App() {
  const [selectedFile, setSelectedFile] = useState(null);
  const [status, setStatus] = useState('');
  const [isUploading, setIsUploading] = useState(false);
  const [documents, setDocuments] = useState([]);
  const [editingId, setEditingId] = useState(null);
  const [editingName, setEditingName] = useState('');

  const loadDocuments = async () => {
    const response = await fetch(`${apiBaseUrl}/documents`);
    const data = await response.json();
    setDocuments(Array.isArray(data) ? data : []);
  };

  useEffect(() => {
    loadDocuments().catch(() => {
      setStatus('Could not load documents.');
    });
  }, []);

  useEffect(() => {
    const intervalId = setInterval(() => {
      loadDocuments().catch(() => {
        // Keep polling quiet; the main status area is enough for user actions.
      });
    }, 4000);

    return () => clearInterval(intervalId);
  }, []);

  const handleFileChange = (event) => {
    const file = event.target.files?.[0] ?? null;
    setSelectedFile(file);
    setStatus('');
  };

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (!selectedFile) {
      setStatus('Choose a PDF or DOCX file first.');
      return;
    }

    const formData = new FormData();
    formData.append('file', selectedFile);

    try {
      setIsUploading(true);
      setStatus('Uploading file...');

      const response = await fetch(`${apiBaseUrl}/documents`, {
        method: 'POST',
        body: formData,
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Upload failed.');
      }

      setStatus(`Uploaded: ${data.fileName}`);
      setSelectedFile(null);
      await loadDocuments();
    } catch (error) {
      setStatus(error.message);
    } finally {
      setIsUploading(false);
    }
  };

  const handleDelete = async (id) => {
    const response = await fetch(`${apiBaseUrl}/documents/${id}`, {
      method: 'DELETE',
    });
    const data = await response.json();

    if (!response.ok) {
      setStatus(data.message || 'Delete failed.');
      return;
    }

    setStatus('Document deleted.');
    await loadDocuments();
  };

  const handleUpdate = async (id) => {
    const response = await fetch(`${apiBaseUrl}/documents/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ fileName: editingName }),
    });
    const data = await response.json();

    if (!response.ok) {
      setStatus(data.message || 'Update failed.');
      return;
    }

    setStatus(`Updated: ${data.fileName}`);
    setEditingId(null);
    setEditingName('');
    await loadDocuments();
  };

  return (
    <main className="app">
      <section className="card">
        <p className="eyebrow">Doc Hub</p>
        <h1>Upload a PDF or DOCX</h1>
        <p className="description">
          Pick a file and send it to your Go backend. The server will accept only
          PDF and DOCX documents.
        </p>

        <form className="upload-form" onSubmit={handleSubmit}>
          <label className="file-picker">
            <span>Select document</span>
            <input
              type="file"
              accept=".pdf,.docx,application/pdf,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
              onChange={handleFileChange}
            />
          </label>

          <button type="submit" disabled={isUploading}>
            {isUploading ? 'Uploading...' : 'Send to API'}
          </button>
        </form>

        <div className="file-meta">
          {selectedFile ? `Selected: ${selectedFile.name}` : 'No file selected yet.'}
        </div>

        {status ? <div className="status">{status}</div> : null}

        <section className="documents-section">
          <div className="section-head">
            <h2>Documents</h2>
            <button type="button" className="secondary-button" onClick={() => loadDocuments()}>
              Refresh
            </button>
          </div>

          <div className="documents-list">
            {documents.length === 0 ? (
              <div className="empty-state">No documents uploaded yet.</div>
            ) : (
              documents.map((document) => (
                <article key={document.id} className="document-row">
                  <div className="document-main">
                    {editingId === document.id ? (
                      <input
                        className="text-input"
                        value={editingName}
                        onChange={(event) => setEditingName(event.target.value)}
                      />
                    ) : (
                      <strong>{document.fileName}</strong>
                    )}
                    <span>{document.contentType}</span>
                    <span>{Math.round(document.size / 1024)} KB</span>
                    <span>Status: {document.status}</span>
                  </div>

                  <div className="document-actions">
                    {editingId === document.id ? (
                      <button type="button" onClick={() => handleUpdate(document.id)}>
                        Save
                      </button>
                    ) : (
                      <button
                        type="button"
                        className="secondary-button"
                        onClick={() => {
                          setEditingId(document.id);
                          setEditingName(document.fileName);
                        }}
                      >
                        Rename
                      </button>
                    )}

                    <button
                      type="button"
                      className="danger-button"
                      onClick={() => handleDelete(document.id)}
                    >
                      Delete
                    </button>
                  </div>
                </article>
              ))
            )}
          </div>
        </section>
      </section>
    </main>
  );
}

export default App;
