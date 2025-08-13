const { chromium } = require('playwright');
const express = require('express');
const { createServer } = require('http');
const { Server } = require('ws');

const app = express();
const server = createServer(app);
const wss = new Server({ server });

wss.on('connection', async (ws) => {
    console.log('Client connected');
    let browser;
    try {
        browser = await chromium.launch();
    } catch (error) {
        console.error('Failed to launch browser:', error);
        ws.send(JSON.stringify({ status: 'error', message: 'Failed to launch browser' }));
        ws.close();
        return;
    }

    ws.on('message', async (message) => {
        try {
            console.log('Received message:', message.toString());
            const { action, url, path } = JSON.parse(message);
            if (action === 'screenshot') {
                if (!url || !path) {
                    throw new Error('Missing url or path in request');
                }
                console.log(`Processing screenshot for URL: ${url}, Path: ${path}`);
                const page = await browser.newPage();
                await page.goto(url, { waitUntil: 'networkidle' });
                await page.screenshot({ path });
                await page.close();
                ws.send(JSON.stringify({ status: 'success', path }));
                console.log('Screenshot saved successfully');
            } else {
                throw new Error('Unknown action');
            }
        } catch (error) {
            console.error('Error processing message:', error);
            ws.send(JSON.stringify({ status: 'error', message: error.message }));
        }
    });

    ws.on('error', (error) => {
        console.error('WebSocket error:', error);
    });

    ws.on('close', async () => {
        console.log('Client disconnected');
        if (browser) {
            await browser.close();
        }
    });
});

server.listen(3000, () => console.log('Playwright server running on ws://localhost:3000'));
