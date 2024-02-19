import { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from 'axios';
import { apiHost } from "./config";

const FileUpload = () => {
    const [file, setFile] = useState();
    const [description, setDescription] = useState('');
    const [uploadedID, setUploadedID] = useState(null);
    const navigate = useNavigate();

    const uploadURL = `${apiHost}/api/files/upload`;

    const handleSubmit = (event) => {
        event.preventDefault()
        setUploadedID(null)
        const formData = new FormData();
        formData.append('upload_file', file);
        formData.append('description', description);

        const config = {
            headers: {
                'content-type': 'multipart/form-data',
            },
        };
        axios.post(uploadURL, formData, config).then((response) => {
            console.log(response.data);
            setUploadedID(response.data.id);

            setTimeout(() => navigate("/files"), 2 * 1000);
        })
            .catch(error => {
                setDescription('');
                setFile();
                setUploadedID(null);

                console.error(error);
            });

    }

    return (
        <div>
            <h2 className="heading">Upload File</h2>
            {uploadedID &&
                <p>Your file has been successfully uploaded. Redirecting to the listings page in 2 seconds.</p>
            }
            <form onSubmit={handleSubmit} className="upload-form">
                <label for="file" className="form-fields-left">File</label>
                <input type="file" onChange={(event) => { setFile(event.target.files[0]) }} className="form-fields-right" />

                <label for="description" className="form-fields-left">Description</label>
                <input type="text" value={description} onChange={(event) => { setDescription(event.target.value) }} className="form-fields-right" />
                <br />
                <button type="submit" className="upload-btn">Upload</button>
            </form>

        </div>
    );
}

export default FileUpload;