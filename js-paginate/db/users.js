// @ts-check
const sqlite3 = require("sqlite3").verbose();
const { Table } = require("console-table-printer");
const { Database } = require("sqlite3");

/**
 * @readonly
 * @enum {Error}
 */
const dbErrors = {
    /** @readonly */
    ErrCouldNotOpen: new Error("Could not open db file"),

    /** @readonly */
    ErrCouldNotCreate: new Error("Could not create table"),

    /** @readonly */
    ErrInsertFailed: new Error("Could not insert user"),

    /** @readonly */
    ErrRetrieveFailed: new Error("Could not retrieve users"),

    /** @readonly */
    ErrInvalidOrder: new Error("Invalid order direction"),

    /** @readonly */
    ErrMissingParameter: new Error("Missing parameter"),
};

class UsersTable {
    /**
     * @param {Database} db
     */
    constructor(db) {
        /** @type {Database} */
        this.db = db;
    }

    static async create() {
        const db = new sqlite3.Database("./users.db", (err) => {
            if (err) {
                throw dbErrors.ErrCouldNotOpen;
            }
        });

        const createTableQuery = `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            age INTEGER NOT NULL,
            country TEXT NOT NULL,
            degree TEXT,
            status TEXT,
            site TEXT
        );
        `;

        return new Promise((res, rej) => {
            db.run(createTableQuery, (err) => {
                if (err) {
                    rej(dbErrors.ErrCouldNotCreate);
                } else {
                    res(new UsersTable(db));
                }
            });
        });
    }

    /** @param {User} user */
    insert(user) {
        const insertQuery = `
        INSERT INTO users (name, age, country, degree, status, site)
        VALUES (?, ?, ?, ?, ?, ?);
        `;

        return new Promise((resolve, reject) => {
            this.db.run(
                insertQuery,
                [
                    user.name,
                    user.age,
                    user.country,
                    user.degree,
                    user.status,
                    user.site,
                ],
                function (err) {
                    if (err) {
                        reject(dbErrors.ErrInsertFailed);
                    }
                },
            );
        });
    }

    /**
     * @param {User[]} users
     * @returns {Promise<string>}
     */
    async show(users) {
        const table = new Table({
            columns: [
                { name: "ID", alignment: "right" },
                { name: "Name", alignment: "left" },
                { name: "Age", alignment: "right" },
                { name: "Country", alignment: "left" },
                { name: "Degree", alignment: "left" },
                { name: "Status", alignment: "left" },
                { name: "Site", alignment: "left" },
            ],
            shouldDisableColors: true,
        });

        users.forEach((user) => {
            table.addRow({
                ID: user.id,
                Name: user.name,
                Age: user.age,
                Country: user.country,
                Degree: user.degree,
                Status: user.status,
                Site: user.site,
            });
        });

        return table.render();
    }

    /**
     * @param {Filters} [filters={}]
     * @param {RetrieveOptions} [options={}]
     */
    retrieve(filters = {}, options = {}) {
        let query = `SELECT id, name, age, country, degree, status, site FROM users WHERE 1=1`;
        const params = [];

        if (filters.status) {
            query += ` AND LOWER(status) = LOWER(?)`;
            params.push(filters.status);
        }

        if (filters.countries && filters.countries.length > 0) {
            query += ` AND LOWER(country) IN (${filters.countries.map(() => "LOWER(?)").join(", ")})`;
            params.push(...filters.countries);
        }

        if (filters.age) {
            query += ` AND age = ?`;
            params.push(filters.age);
        }

        if (filters.degree) {
            query += ` AND LOWER(degree) = LOWER(?)`;
            params.push(filters.degree);
        }

        if (options.orderBy && options.order) {
            if (options.order !== "asc" && options.order !== "desc") {
                return reject(ErrInvalidOrder);
            }
            query += ` ORDER BY ${options.orderBy} ${options.order}`;
        }

        if (options.limit > 0) {
            query += ` LIMIT ?`;
            params.push(options.limit);
        }

        return new Promise((resolve, reject) => {
            this.db.all(query, params, (err, rows) => {
                if (err) {
                    reject(dbErrors.ErrRetrieveFailed);
                } else {
                    resolve(
                        rows.map((row) => ({
                            id: row.id,
                            name: row.name,
                            age: row.age,
                            country: row.country,
                            degree: row.degree,
                            status: row.status,
                            site: row.site,
                        })),
                    );
                }
            });
        });
    }

    /**
     * @param {number} pageSize
     * @param {number} page
     * @param {OrderBy} orderBy
     * @param {Filters} filters
     * @returns {Promise<User[]>}
     */
    paginate(pageSize, page, orderBy, filters) {
        let query = `SELECT id, name, age, country, degree, status, site FROM users WHERE 1=1`;
        const params = [];

        query = this.#applyFilters(query, params, filters);

        query = this.#applyOrderBy(query, orderBy);

        const offset = (page - 1) * pageSize;
        query += ` LIMIT ? OFFSET ?`;
        params.push(pageSize, offset);

        return new Promise((resolve, reject) => {
            this.db.all(query, params, (err, rows) => {
                if (err) {
                    return reject(err);
                }
                resolve(
                    rows.map((row) => ({
                        id: row.id,
                        name: row.name,
                        age: row.age,
                        country: row.country,
                        degree: row.degree,
                        status: row.status,
                        site: row.site,
                    })),
                );
            });
        });
    }

    close() {
        this.db.close();
    }

    /**
     * @param {string} query
     * @param {any[]} params
     * @param {Filters} filters
     * @returns {string}
     */
    #applyFilters(query, params, filters) {
        if (!filters) {
            return query;
        }

        if (filters.status) {
            query += ` AND LOWER(status) = LOWER(?)`;
            params.push(filters.status);
        }

        if (filters.countries && filters.countries.length > 0) {
            query += ` AND LOWER(country) IN (${filters.countries.map(() => "LOWER(?)").join(", ")})`;
            params.push(...filters.countries);
        }

        if (filters.age) {
            query += ` AND age = ?`;
            params.push(filters.age);
        }

        if (filters.degree) {
            query += ` AND LOWER(degree) = LOWER(?)`;
            params.push(filters.degree);
        }

        return query;
    }

    /**
     * @param {string} query
     * @param {OrderBy} orderBy
     * @returns {string}
     */
    #applyOrderBy(query, orderBy) {
        if (!orderBy) {
            return query;
        }
        /** @type {string[]} */
        const orderClauses = [];
        if (orderBy.id) orderClauses.push(`id ${orderBy.id}`);
        if (orderBy.age) orderClauses.push(`age ${orderBy.age}`);
        if (orderBy.name) orderClauses.push(`name ${orderBy.name}`);
        if (orderClauses.length > 0) {
            query += ` ORDER BY ${orderClauses.join(", ")}`;
        }

        return query;
    }
}

module.exports = {
    dbErrors,
    UsersTable,
};
