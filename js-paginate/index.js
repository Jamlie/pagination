// @ts-check
const express = require("express");
const usersRoutes = require("./routes/pagination");
const { initializeUsersTable } = require("./model/users");

/** @typedef {import("express").Request} ExpressRequest */
/** @typedef {import("express").Response} ExpressResponse */

const app = express();
const port = process.env.PORT || 8080;

app.use(express.json());

initializeUsersTable();

app.use("/", usersRoutes);

app.listen(port, () => {
    console.log(`Server running on http://localhost:${port}`);
});
