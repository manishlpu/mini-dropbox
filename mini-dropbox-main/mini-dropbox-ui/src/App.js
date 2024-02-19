import { Routes, Route, Link } from "react-router-dom";
import Files from "./components/Files";
import FileUpload from "./components/FileUpload";
import File from "./components/File";
import FilesPage from "./components/FilesPage";
import NoMatch from "./components/NoMatch";
import './App.css';

const App = () => {
  return (
    <>
      <nav className="navbar">
        <div className="navbar-brand">Mini Dropbox</div>
        <ul className="navbar-list">
          <li className="navbar-item">
            <Link to="/" className="navbar-item-link">Home</Link>
          </li>
          <li className="navbar-item">
            <Link to="/files" className="navbar-item-link">Files</Link>
          </li>
        </ul>
      </nav>

      <Routes>
        <Route path="/" element={<FileUpload />} />
        <Route path="/files" element={<FilesPage />}>
          <Route index element={<Files />} />
          <Route path=":fileId" element={<File />} />
        </Route>
        <Route path="*" element={<NoMatch />} />
      </Routes>
    </>
  );
};

export default App;
