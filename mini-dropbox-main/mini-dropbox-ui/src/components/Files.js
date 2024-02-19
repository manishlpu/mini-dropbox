import { useState, useEffect } from 'react';
import { Link } from "react-router-dom";
import { apiHost } from './config';
import axios from 'axios';

const Files = () => {
    const [files, setFiles] = useState([]);
    const [loading, setLoading] = useState(true);

    const filesURL = `${apiHost}/api/files`;

    useEffect(() => {
        setLoading(true);

        axios.get(filesURL)
            .then((response) => {            
                setFiles(response.data);
                setLoading(false);
            })
            .catch((error) => {
                console.error('Error fetching data:', error);
                setLoading(false);
            });
    }, []);

    if (loading) {
        return <div>Loading...</div>;
    }

    return (
        <>
            <div className="grid-container">
                {files && files.map((item) => (
                    <div key={item.id} className="grid-item">
                        <img className='file-icon' src='./images/file-logo.png' alt='file' />
                        <span className='file-name'>{item.filename}</span>
                        <Link to={`/files/${item.id}`} className='file-view-btn'>
                            View File
                        </Link>
                    </div>
                ))}
            </div>
        </>
    );
}

export default Files;