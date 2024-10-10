declare type User = {
    id?: number;
    name: string;
    age: number;
    country: string;
    degree?: string;
    status?: string;
    site?: string;
};

declare type OrderBy = {
    id?: string;
    age?: string;
    name?: string;
};

declare type Filters = {
    status?: string;
    age?: number;
    degree?: string;
    countries?: string[];
};

declare type RetrieveOptions = {
    orderBy?: string;
    order?: string;
    limit?: number;
};

declare type APIResponse = {
    data?: User[];
    message?: string;
    statusCode: number;
};
