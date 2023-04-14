const express = require("express");
const crypto = require("node:crypto");
const { readFileSync } = require("node:fs");
const { resolve } = require("node:path");

const app = express();
const port = 3000;

app.use(express.json());

app.get("/", (req, res) => {
  res.send("Hello World!");
});

app.post("/", async (req, res) => {
  const message = req.body.message;
  if (!message) res.status(400).send();

  let pemContents = readFileSync(
    resolve(__dirname, "../public-key.pem"),
    "utf-8"
  );
  pemContents = pemContents.replace("-----BEGIN PUBLIC KEY-----\n", "");
  pemContents = pemContents.replace("-----END PUBLIC KEY-----\n", "");
  const binaryDer = Buffer.from(pemContents, "base64");
  const key = await crypto.subtle.importKey(
    "spki",
    binaryDer,
    {
      name: "RSA-OAEP",
      hash: "SHA-256",
    },
    true,
    ["encrypt"]
  );

  const encrypted = await crypto.subtle.encrypt(
    "RSA-OAEP",
    key,
    new TextEncoder().encode(message)
  );

  console.log("sending", message);
  const encryptedStr = Buffer.from(encrypted).toString("base64");
  console.log("encrypted", encryptedStr);
  const result = await fetch("http://localhost:1323", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ message: encryptedStr }),
  });
  const response = await result.json();
  console.log(response);
  res.send(response);
});

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`);
});
