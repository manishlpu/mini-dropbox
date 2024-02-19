import { useState, useEffect } from 'react';
import { useParams, useNavigate } from "react-router-dom";
import axios from 'axios';
import { apiHost } from './config';

const File = () => {
    const { fileId } = useParams();
    const [file, setFile] = useState({});
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(false);
    const [action, setAction] = useState('');

    const [name, setName] = useState();
    const [description, setDescription] = useState('');
    const [uploadedID, setUploadedID] = useState(null);
    const navigate = useNavigate();

    const fileURL = `${apiHost}/api/files/${fileId}`;

    const updateFile = (event) => {
        event.preventDefault();

        const formData = new FormData();
        formData.append('upload_file', name);
        formData.append('description', description);

        const config = {
            headers: {
                'content-type': 'multipart/form-data',
            },
        };
        axios.putForm(fileURL, formData, config).then((response) => {
            console.log(response.data);
            setUploadedID(response.data.id);
            setAction('updated');

            setTimeout(() => navigate("/files"), 2 * 1000);
        })
            .catch((error) => {
                setDescription('');
                setName();
                setUploadedID(null);

                console.error(error);
            });
    }

    const deleteFile = () => {
        // event.preventDefault();
        setLoading(true);

        axios.delete(fileURL).then((response) => {
            console.log(response.data);
            setAction('deleted');

            setTimeout(() => navigate("/files"), 2 * 1000);
        }).catch((error) => {
            setError(true);
            console.error(error);
        });

        setLoading(false);
    }

    useEffect(() => {
        setLoading(true);

        axios.get(fileURL)
            .then((response) => {
                setFile(response.data);
                setError(false);
            })
            .catch((error) => {
                console.error('Error fetching data:', error);
                setError(true);
            });

        setLoading(false);
    }, []);


    if (loading) {
        return <div>Loading...</div>;
    }

    return (
        <>
            <h2>Detailed File View</h2>
            {!loading && error && <p>Error fetching data.</p>}

            {action && (<p>Your file is {action}! You will be redirected to files listing in a few seconds. Please do not refresh.</p>)
            }
            {!loading && !error &&
                <article className='file-view'>
                    <section className='left'>
                        <img className='file-icon-big' src='./images/file-logo.png' alt='file' />
                        <p className='filename'>{file.filename}</p>

                        <div className='download btn'>
                            <a href={file.s3_object_key} className='btn-link'>Download</a>
                        </div>
                    </section>
                    <section className='right'>
                        {file.description == '' ? '' : (
                            <p className='description'>
                                <span className='title'>Description:</span>
                                <span className='data'>{file.description}</span>
                            </p>
                        )
                        }

                        <p className='size'>
                            <span className='title'>Size:</span>
                            <span className='data size'>{file.size_in_bytes > 0 ? file.size_in_bytes + ' bytes' : 'unknown'}</span>
                        </p>

                        <p className='created-at'>
                            <span className='title'>Created At:</span>
                            <span className='data created-at'>{file.created_at}</span>
                        </p>

                        <p className='updated-at'>
                            <span className='title'>Last Updated:</span>
                            <span className='data updated-at'>{file.updated_at}</span>
                        </p>

                        <article className='adjacent-btns'>
                            <div className='btn update'>
                                <a href='#' onClick={updateFile} className='btn-link'>Update</a>
                            </div>
                            <div className='btn delete'>
                                <a href='#' onClick={deleteFile} className='btn-link'>Delete</a>
                            </div>
                        </article>
                        <form className="upload-form">
                            <label for="file" className="form-fields-left">File</label>
                            <input type="file" onChange={(event) => { setName(event.target.files[0]) }} className="form-fields-right" />

                            <label for="description" className="form-fields-left">Description</label>
                            <input type="text" value={description} onChange={(event) => { setDescription(event.target.value) }} className="form-fields-right" />
                            <br />
                        </form>
                    </section>
                </article>
            }
        </>
    );
}

export default File;