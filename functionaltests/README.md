# CRB Functional Tests #
# Pre-requisites:

 * [Docker](https://docs.docker.com/engine/installation/) is installed and running.
 * CRB is deployed and running on PCF.

# Usage:
The tests are executed using a [container](Dockerfile) which has the required test tools and the copy repository configured.

  * Build the docker image ocopea_crb_test_tool from the functionaltests folder. The image contains all the test tools, test scripts as well as acts as a copy repository for the tests.
  Note: Feel free to change names and tags, and change the following commands accordingly
  ```
  docker build -t ocopea_crb_test_tool .
  ```

  * Find an unused port on the local host that will be used to forward traffic to the container's ssh port i.e. port 22. Ocopea's CRB service uses SFTP to securely transfer copies on this port.
  
  On windows or linux based machines to get the list of *used* local port, run
  ```
  netstat -ap
  ```
  and accordingly choose an unused port for SSH forwarding. This port will be used in later commands.

  * Run the docker container with name ocopea_crb_test_tool

    * If using existing test scripts:
    ```
    docker run -d -p <ssh_fwd_port>:22 --name ocopea_crb_test_tool ocopea_crb_test_tool
    ```

    * If testing new test scripts please use other docker options such as rebuilding the image, mounting the test folder or copy the new files to a running image. Here is an example on linux with mounted volume:
    ```
    docker run -d -p <ssh_fwd_port>:22 -v $(pwd):/root/functionaltests/ --name ocopea_crb_test_tool ocopea_crb_test_tool
    ```

  * Run the tests using the above container
  ```
  docker exec ocopea_crb_test_tool jasmine-node main_test_spec.js --config host "<crb_app_url>" --config dbname "<p-mysql_name>" --config dbpassword "<p-mysql_password>" --config dbuser "<p-mysql_username>" --config dbip "<p-mysql_hostname>" --config repoip "<host_address>" --config repoport "<ssh_fwd_port>" --config repouser "root" --config repopassword "screencast" --config crbname "<crb_app_name>" --config cfusername "<pcf_username>" --config cfpassword "<pcf_password>" --config cfendpoint "<pcf_api_endpoint>" --config version "<crb_app_version>"
  ```
  **Note**:

  * <crb_app_name> and <crb_app_url> are the name and url of the crb running in pcf. They can be retrieved by running
  ```
  cf a
  ```
  Here is a sample output
  ```
  name         requested state   instances   memory   disk   urls
  crb-server   started           1/1         1G       1G     crb-server.cf.isus.emc.com
  ```

  * <p-mysql_name>, <p-mysql_username>, <p-mysql_password> and <p-mysql_hostname> are the credentials of the p-mysql service which is bound to the crb app. They can be retrieved by running
  ```
  cf env <crb_app_name>
  ```

  Here is an sample output:
  ```
  System-Provided:
  {
   "VCAP_SERVICES": {
    "p-mysql": [
     {
      "credentials": {
      "hostname": "10.106.124.194",
      "jdbcUrl": "jdbc:mysql://10.106.124.194:3306/cf_b5ffc129_539c_42cd_ac9b_5524f93502c9?user=njT6gHoQZrca1Yzw\u0026password=8FIWMg5otqKpWKvM",
      "name": "cf_b5ffc129_539c_42cd_ac9b_5524f93502c9",
      "password": "8FIWMg5otqKpWKvM",
      "port": 3306,
      "uri": "mysql://njT6gHoQZrca1Yzw:8FIWMg5otqKpWKvM@10.106.124.194:3306/cf_b5ffc129_539c_42cd_ac9b_5524f93502c9? reconnect=true",
      "username": "njT6gHoQZrca1Yzw"
     },
     "label": "p-mysql",
     "name": "crb-mysql",
  ...
  }
  ... 
  ```
  * <host_address> is the ip of the copy repo system.
  * <ssh_fwd_port> is the port used for SSH forwarding.
  * <pcf_api_endpoint> is the API end point of the PCF. This can be retrieved by running
  ```
  cf t
  ```
  Note: Do not include https:\\ or http:\\ in <pcf_api_endpoint>
  * <pcf_username>, <pcf_password> are the credentials for pcf end point.
  * <crb_app_version> is the version of the CRB. This can be retreived from the [manifest](../manifest.yml).

  #### test output
  A successfull test run should give an output as shown below
  ```
  .testDeleteCopy : Posted Repo
  .testDeleteCopy : Posted Copy
  .testDeleteCopy : Deleted Copy
  .
  .
  .
  .testPostCopyWithNoDBAccess :Get copy with no DB access
  cleaning up db connection
  .

  Finished in 43.48 seconds
  86 tests, 209 assertions, 0 failures, 0 skipped
  ```