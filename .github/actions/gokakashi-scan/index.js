import fetch from 'node-fetch';
import { setFailed, getInput, setOutput } from '@actions/core';

const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms));

async function makeRequest(url, options) {
    try {
        const response = await fetch(url, options);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error(`Request failed: ${error.message}`);
        throw error;
    }
}

async function pollScanStatus(apiHost, apiToken, scanId) {
    let status = 'queued';
    let retries = 0;
    const maxRetries = 30;
    const initialDelay = 10000;

    while ((status === 'queued' || status === 'in-progress') && retries < maxRetries) {
        console.log(`Current scan status: ${status}. Waiting for completion...`);

        await sleep(initialDelay * Math.pow(2, retries));

        const statusData = await makeRequest(`${apiHost}/api/v0/scan/${scanId}/status`, {
            method: 'GET',
            headers: { 'Authorization': `Bearer ${apiToken}` }
        });
        console.log(statusData)
        status = statusData.status;
        retries++;

        if (status === 'completed') {
            return statusData.report_url[0];
        }
    }

    throw new Error(`Scan did not complete within the expected time. Final status: ${status}`);
}

async function run() {
    try {
        const apiHost = getInput('api_host');
        const apiToken = getInput('api_token');
        const imageName = getInput('image_name');
        const severity = getInput('severity');
        const publish = getInput('publish');
        const failOnSeverity = getInput('fail_on_severity');
        console.log(`API Host: ${apiHost}`);
        console.log(`API Token: ${apiToken}`);
        console.log(`Image Name: ${imageName}`);

        // Validate inputs
        if (!apiHost || !apiToken || !imageName) {
            throw new Error('Missing required inputs: api_host, api_token, or image_name');
        }

        console.log('Triggering scan...');
        const triggerData = await makeRequest(`${apiHost}/api/v0/scan?image=${imageName}&severity=${severity}&publish=${publish}`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${apiToken}`,
                'Content-Type': 'application/json'
            }
        });

        const scanId = triggerData.scan_id;
        console.log(`Scan triggered with scan ID: ${scanId}`);

        const reportUrl = await pollScanStatus(apiHost, apiToken, scanId);
        console.log(`Scan completed. Report URL: ${reportUrl}`);
        setOutput('report_url', reportUrl);

        console.log('Fetching scan report...');
        const reportData = await makeRequest(reportUrl);

        if (failOnSeverity) {
            const severitiesToFailOn = failOnSeverity.split(',').map(sev => sev.trim().toUpperCase());
            const hasVulnsToFail = reportData.vulnerabilities.some(vuln =>
                severitiesToFailOn.includes(vuln.severity)
            );
            if (hasVulnsToFail) {
                setFailed(`Vulnerabilities found with severity: ${severitiesToFailOn.join(', ')}`);
            } else {
                console.log(`No vulnerabilities found with severity: ${severitiesToFailOn.join(', ')}`);
            }
        } else {
            console.log('No fail_on_severity defined, proceeding without failing the job.');
        }
    } catch (error) {
        setFailed(`Action failed: ${error.message}`);
    }
}

run();