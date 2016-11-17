
const URL_PREFIX = '/api/v1';

export interface QueryResult<T> {
    from?: number,
    to?: number,
    count?: number,
    result: T[];
}

export default class WebApi {
    static get<T>(path: string): Promise<T> {
        return fetch(toUrl(path))
            .then(this.toJson);
    }

    static post<T>(path: string, jsonBody: any): Promise<T> {
        return fetch(toUrl(path), {
            method: 'post',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(jsonBody)
        })
            .then(this.toJson);
    }

    static put<T>(path: string, _rev: string = '*', jsonBody: any): Promise<T> {
        return fetch(toUrl(path), {
            method: 'put',
            headers: {
                'If-Match': _rev,
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(jsonBody)
        })
            .then(this.toJson);
    }

    static del<T>(path: string, _rev: string = '*'): Promise<T> {
        return fetch(toUrl(path), {
            method: 'delete',
            headers: {
                'If-Match': _rev
            }
        })
            .then(this.toJson);
    }

    static patch<T>(path: string, _rev: string, jsonBody: Object): Promise<T> {
        return fetch(toUrl(path), {
            method: 'patch',
            headers: {
                'If-Match': _rev,
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(jsonBody)
        })
            .then(this.toJson);
    }

    // CRUD Query
    static create<T>(resource: string, _id: string, jsonBody: Object): Promise<T> {
        if (_id) {
            jsonBody['_id'] = _id;
        }
        return this.post(resource, jsonBody)
    }

    static read<T>(fullId: string): Promise<T> {
        return this.get(fullId)
    }

    static update<T>(fullId: string, _rev: string, jsonBody: Object): Promise<T> {
        return this.put(fullId, _rev, jsonBody)
    }

    static delete<T>(fullId: string, _rev: string): Promise<T> {
        return this.del(fullId, _rev)
    }

    static query<T>(resource: string, queryParams: Object | string = {}): Promise<QueryResult<T>> {
        let q = queryParams;
        if (typeof queryParams !== 'string') {
            q = this.queryString(queryParams);
        }
        return fetch(toUrl(`${resource}?${q}`))
            .then(this.toJson);
    }

    static action<T>(resource: string, actionId: string, values?: Object): Promise<T> {
        return this.post(`${resource}?_action=${actionId}`, values);
    }

    // Utils
    static toText(response) {
        return response.text();
    }

    static toJson(response) {
        return response.json();
    }

    static queryString(queryParams: Object = {}): string {
        const query = Object.keys(queryParams).map(x => {
            const v = queryParams[x];
            if (Array.isArray(v)) {
                return v.map(y => {
                    return `${x}=${y}`;
                }).join('&');
            } else {
                return `${x}=${v}`;
            }
        }).join('&');
        return query;
    }
}

function toUrl(path: string) {
    return `${URL_PREFIX}/${path}`;
}