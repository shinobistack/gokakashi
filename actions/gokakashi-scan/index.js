import fetch from 'node-fetch';
import { setFailed, getInput, setOutput } from '@actions/core';

async function run() {
    try {
        // Input parameters from the action
        const apiHost = getInput('api_host');
        const apiToken = getInput('api_token');
        const imageName = getInput('image_name');
        const severity = getInput('severity');
        const publish = getInput('publish');

        // Step 1: Trigger the scan and get the scan_id
        const triggerResponse = await fetch(`${apiHost}/api/v0/scan?image=${imageName}&severity=${severity}&publish=${publish}`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${apiToken}`,
                'Content-Type': 'application/json'
            }
        });

        if (!triggerResponse.ok) {
            throw new Error(`Failed to trigger the scan. Status: ${triggerResponse.status}`);
        }

        const triggerData = await triggerResponse.json();
        const scanId = triggerData.scan_id;

        console.log(`Scan triggered with scan ID: ${scanId}`);

        // Step 2: Poll the scan status until it's completed
        let status = 'queued';
        let reportUrl = '';

        while (status === 'queued' || status === 'in-progress') {
            console.log(`Current scan status: ${status}. Waiting for completion...`);

            await new Promise(r => setTimeout(r, 10000)); // Wait 10 seconds between polls

            const statusResponse = await fetch(`${apiHost}/api/v0/scan/${scanId}/status`, {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${apiToken}`
                }
            });

            if (!statusResponse.ok) {
                throw new Error(`Failed to get scan status. Status: ${statusResponse.status}`);
            }

            const statusData = await statusResponse.json();
            status = statusData.status;

            // Check if scan is completed
            if (status === 'completed') {
                reportUrl = statusData.report_url[0]; // Extract report URL
                console.log(`Scan completed. Report URL: ${reportUrl}`);
                setOutput('report_url', reportUrl);  // Set the output for future steps
            }
        }

        // If the scan did not complete successfully
        if (status !== 'completed') {
            throw new Error(`Scan failed with status: ${status}`);
        }

        // Step 3: Check the scan report for vulnerabilities
        const reportResponse = await fetch(reportUrl);
        const reportData = await reportResponse.json();

        const hasCriticalVulns = reportData.vulnerabilities.some(vuln => vuln.severity === 'CRITICAL');

        if (hasCriticalVulns) {
            setFailed('Critical vulnerabilities found in the image.');
        } else {
            console.log('No critical vulnerabilities found.');
        }

    } catch (error) {
        setFailed(error.message);
    }
}

run();
