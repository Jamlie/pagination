// @ts-check
const { UsersTable } = require("../db/users");

/** @type {UsersTable | null} */
let usersTable = null;

/**
 * Initializes the UsersTable
 */
async function initializeUsersTable() {
    try {
        usersTable = await UsersTable.create();
        return usersTable;
    } catch (error) {
        console.error("Failed to initialize users table", error);
        process.exit(1);
    }
}

/**
 * Inserts a user into the table
 * @param {Object} user
 */
async function insertUser(user) {
    return usersTable?.insert(user);
}

/**
 * Retrieves all users
 */
async function getAllUsers() {
    return usersTable?.retrieve();
}

/**
 * @param {Filters} [filters={}]
 * @param {RetrieveOptions} [options={}]
 */
async function getUsersWith(filters = {}, options = {}) {
    return usersTable?.retrieve(filters, options);
}

/**
 * Displays users in a formatted way
 * @param {Array} users
 */
async function showUsers(users) {
    return usersTable?.show(users);
}

/**
 * Paginate users
 * @param {number} pageSize
 * @param {number} page
 * @param {OrderBy} orderBy
 * @param {Filters} filters
 */
async function paginateUsers(pageSize, page, orderBy, filters) {
    return usersTable?.paginate(pageSize, page, orderBy, filters);
}

module.exports = {
    initializeUsersTable,
    insertUser,
    getAllUsers,
    getUsersWith,
    showUsers,
    paginateUsers,
};
