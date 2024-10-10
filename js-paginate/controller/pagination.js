// @ts-check
const {
    insertUser,
    getAllUsers,
    showUsers,
    paginateUsers,
} = require("../model/users");
const { dbErrors } = require("../db/users");

/** @typedef {import("express").Request} ExpressRequest */
/** @typedef {import("express").Response} ExpressResponse */

/**
 * Sends a JSON response
 * @param {ExpressResponse} res
 * @param {APIResponse} response
 */
function sendJson(res, response) {
    res.status(response.statusCode).json(response);
}

/**
 * Inserts a user into the database
 * @param {ExpressRequest} req
 * @param {ExpressResponse} res
 */
async function insert(req, res) {
    const user = req.body;

    try {
        await insertUser(user);
        sendJson(res, {
            statusCode: 200,
            message: "User inserted successfully",
        });
    } catch (error) {
        console.error(error);
        sendJson(res, {
            statusCode: 500,
            message: "Failed to insert user",
        });
    }
}

/**
 * Displays all users
 * @param {ExpressRequest} _req
 * @param {ExpressResponse} res
 */
async function show(_req, res) {
    try {
        const allUsers = await getAllUsers();
        const buf = await showUsers(allUsers);

        res.status(200).send(buf);
    } catch (error) {
        console.error(error);
        sendJson(res, {
            statusCode: 500,
            message: "Failed to retrieve users",
        });
    }
}

/**
 * Paginates users
 * @param {ExpressRequest} req
 * @param {ExpressResponse} res
 */
async function paginate(req, res) {
    const { pageSize, page, orderBy, filters } = req.body;

    if (!pageSize || !page) {
        return sendJson(res, {
            statusCode: 400,
            message: dbErrors.ErrMissingParameter.message,
        });
    }

    try {
        const retrievedUsers = await paginateUsers(
            pageSize,
            page,
            orderBy,
            filters,
        );
        sendJson(res, {
            statusCode: 200,
            data: retrievedUsers,
        });
    } catch (error) {
        console.error(error);
        sendJson(res, {
            statusCode: 500,
            message: dbErrors.ErrMissingParameter.message,
        });
    }
}

module.exports = {
    insert,
    show,
    paginate,
};
