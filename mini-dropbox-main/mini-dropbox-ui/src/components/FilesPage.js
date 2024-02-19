import { Outlet } from "react-router-dom";

const FilesPage = () => {
    return (

        <section>
            <h2 className='heading'>User Files</h2>

            <Outlet />
        </section>
    );
}

export default FilesPage;