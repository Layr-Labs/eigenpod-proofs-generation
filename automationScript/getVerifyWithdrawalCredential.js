require('dotenv').config();

const { execFile } = require("child_process");
const axios = require('axios');
const fs = require('fs');

const apiKey = '31303734|466e2867dcf14abcbb9513eed794d896';

const headers = {
    'X-API-Key': apiKey,
};

const outputPathForStateFile = 'STATE_FILE.json';
const outputPathForHeadFile = 'HEAD_FILE.json';
const slotNumber = 9179815; //Update with current timeStamp
//BeaconOracle get the latest slot number
//https://etherscan.io/address/0x343907185b71adf0eba9567538314396aa985442
//If doing for different validator, Modify validator in the withdrawalCredential.sh 

const headURL = `https://data.spiceai.io/eth/beacon/eth/v1/beacon/headers/${slotNumber}`;
const stateURL = `https://data.spiceai.io/eth/beacon/eth/v2/debug/beacon/states/${slotNumber}`;

const getHeadFile = async (url, outputPath, headers) => {
    try {
        const response = await axios.get(url, { headers });
        fs.writeFileSync(outputPath, JSON.stringify(response.data, null, 2));
        console.log('File saved successfully.');
    } catch (error) {
        console.error('Error fetching the data:', error);
    }
};

const getStateFile = async (url, outputPath, headers) => {
    try {
        const response = await axios({
            method: 'GET',
            url: url,
            responseType: 'stream',
            headers: headers,
            timeout: 30000 // Adjust timeout as needed
        });

        const writer = fs.createWriteStream(outputPath);

        response.data.pipe(writer);

        return new Promise((resolve, reject) => {
            writer.on('finish', resolve);
            writer.on('error', reject);
        });
    } catch (error) {
        console.error('Error fetching the data:', error);
        throw error;
    }
};

const run = async () => {
    try {
        //Fetch the head file
        await getHeadFile(headURL, outputPathForHeadFile, headers);

        // Fetch the state file
        await getStateFile(stateURL, outputPathForStateFile, headers);
        console.log('Files saved successfully.');

        // Execute the first script
        execFile('../automationScript/intialize.sh', (error, stdout, stderr) => {
            if (error) {
                console.error("Error running intialize.sh:", error);
                return;
            }
            if (stderr) {
                console.error("stderr from intialize.sh:", stderr);
                return;
            }
            console.log("intialize.sh stdout:", stdout);

            // Execute the second script after the first completes successfully
            execFile('../automationScript/withdrawalCredential.sh', (error, stdout, stderr) => {
                if (error) {
                    console.error("Error running withdrawalCredential.sh:", error);
                    return;
                }
                if (stderr) {
                    console.error("stderr from withdrawalCredential.sh:", stderr);
                    return;
                }
                console.log("withdrawalCredential.sh stdout:", stdout);
                console.log("Successfully Created VerifyWithdrawalCredential JSON");
            });
        });
    } catch (error) {
        console.error('Error during operations:', error);
    }
};

run();
