# mini-dropbox

## Description
The goal of this project is to implement a simplified Dropbox-like service where users can upload, retrieve, and manage their files through a set of RESTful APIs. Alongside the backend APIs, a basic UI will also be provided to showcase these functionalities. The service should also support the storage of metadata for each uploaded file, such as the file name, creation timestamp, and more. 

### API Requirements (Functional)
- [X] **POST**    `/files/upload` Allow users to upload files onto the platform.
- [X] **GET**     `/files/{fileID}` Retrieve a specific file based on a unique identifier.
- [X] **PUT**     `/files/{fileID}` Update an existing file or its metadata.
- [X] **DELETE**  `/files/{fileID}` Delete a specific file based on a unique identifier. 
- [X] **GET**     `/file` List all available files and their metadata.

**Note**: Applied a soft delete, instead of hard delete for the file. Wrote a separate cron to delete the file data after 30 days of inactivity.

### User Interface
1. **File Upload Section**: A form to upload a new file and its metadata.
1. **File List Section**: A table or list view that showcases all the files available on the platform.
1. **File Action Section**: Options to download, update, or delete files by interacting with the corresponding APIs.

### Technologies
1. **Backend**: Vanilla Golang with libraries support like *Cobra, Gorilla, SQL ORM*, etc.
2. **Database**: *MySQL* (Relational database) to store the files and metadata.
3. **Frontend**: A basic UI developed using ReactJS Framework.
4. **Storage**: Used *AWS S3* for storing files, which can be configured using environment file.

### Pre-requisites and dependencies
1. Golang v1.19 or above
2. MySQL database (v8.0 or above)
3. make (used to run Makefile)
4. S3 bucket with public access and following server ACLs: GetObject, PutObject, DeleteObject, ListObject, etc.
- Sample S3 bucket Policy: 
    ```json
    {
        "Version": "2012-10-17",
        "Statement": [
            {
                "Sid": "Stmt1405592139000",
                "Effect": "Allow",
                "Principal": "*",
                "Action": "s3:GetObject",
                "Resource": [
                    "arn:aws:s3:::bucketname/*",
                    "arn:aws:s3:::bucketname"
                ]
            }
        ]
    }
    ```

### Steps to run the backend server
1. Navigate to the project directory in terminal and fetch all the dependencies using following command:
    ```sh 
    $ go mod download
    ```
1. Create and update the .env file in the project directory by using `.env.example` file. Fill in all the necessary values to make connection with the pre-requisites defined above.
1. Add the required table(s) to your RDBMS system by using commands from `metadata.sql`. 
**Note**: This is a one-time step only. Taking this step again will clear all the metadata from RDBMS.
1. Run the application using terminal by typing: 
    ```sh
    $ make run
    ```

### Steps to start the UI
1. Navigate to the directory `mini-dropbox-ui` inside the root directory of project. Install all the dependencies using following command: 
    ```sh
    $ npm install
    ```
1. Create and update the .env.local file in the `mini-dropbox-ui` directory with key `REACT_APP_BACKEND_HOST`. Sample: 
    ```sh
    REACT_APP_BACKEND_HOST="http://localhost:8082"
    ```
1. Run the React application on local computer using the following command:
    ```sh
    $ npm start
    ```
1. Open your browser and navigate to the UI using the address: `http://localhost:3000`

### Extra Points considered for Application Server
1. **Graceful shutdown**: This avoids any side effects on conflicts that may occur on closing the server and the new deployment can be started without any kind of difficulty.
1. **Logging**: For debugging and monitoring the application on remote servers, it is recommended to log the application functionality.
1. **Panic Handler**: Used to prevent the application from being killed, in case of any runtime errors or application malfunctioning.

### Improvements that can be done
1. Unit tests
1. Caching the RDBMS response, to save DB queries
1. Dockerfile - to improve collaboration and ease of working
1. Pagination - implement pagination for the listing page to improve performance and efficiency.
