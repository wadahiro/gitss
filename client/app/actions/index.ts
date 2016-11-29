import { Action, Dispatch } from 'redux';
import { browserHistory } from 'react-router';

const { ActionCreators } = require('redux-undo');

import { RootState, FilterParams, Indexed } from '../reducers';
import WebApi from '../api/WebApi';

export type Actions =
    GetIndexedList |
    GetBaseFilters |
    Search |
    ResetFacets |
    SearchStart |
    ToggleSearchOptions
    ;

export interface GetBaseFilters extends Action {
    type: 'GET_BASE_FILTERS';
    payload: {
        organizations: string[];
        projects: string[];
        repositories: string[];
        branches: string[];
        tags: string[];
    };
}

export function getBaseFilters(dispatch: Dispatch<Search>, organization: string, project: string, repository: string): void {
    let urlPath = '';
    if (organization) {
        urlPath += `/${organization}`;
        if (project) {
            urlPath += `/${project}`;
            if (repository) {
                urlPath += `/${repository}`;
            }
        }
    }
    WebApi.get(`filters${urlPath}`)
        .then(res => {
            // console.log(res);
            dispatch({
                type: 'GET_BASE_FILTERS',
                payload: res
            });
        })
        .catch(e => {
            console.warn(e);
        });
}

export interface Search extends Action {
    type: 'SEARCH';
    payload: {
        result: any;
    };
}

export interface ResetFacets extends Action {
    type: 'RESET_FACETS';
}

export interface SearchStart extends Action {
    type: 'SEARCH_START';
    payload: {
        filterParams: FilterParams
    };
}

export function triggerSearch(dispatch: Dispatch<Search>, query?: string): void {
    // reset filters & current facets
    dispatch({
        type: 'RESET_FACETS'
    });

    const params = {
        q: query
    };
    _triggerSearch(params, 0);
}

export function triggerFilter(dispatch: Dispatch<Search>, filterParams: FilterParams, query: string, page: number = 0): void {
    const params = {
        ...filterParams,
        q: query
    };
    _triggerSearch(params, page);
}

export interface ToggleSearchOptions extends Action {
    type: 'TOGGLE_SEARCH_OPTIONS';
}

export function toggleSearchOptions(dispatch: Dispatch<RootState>): void {
    dispatch({
        type: 'TOGGLE_SEARCH_OPTIONS'
    });
}

function _makeQueryString(filterParams: FilterParams, query: string): string {
    const params = {
        ...filterParams,
        q: query
    };
    return Object.keys(params).map(x => {
        return `${x}=${params[x]}`;
    }).join('&');
}

export function search(dispatch: Dispatch<Search>, filterParams: FilterParams): void {

    dispatch({
        type: 'SEARCH_START',
        payload: {
            filterParams
        }
    });

    const params = {
        ...filterParams
    };

    WebApi.query('search', params)
        .then(res => {
            // console.log(res);
            dispatch({
                type: 'SEARCH',
                payload: {
                    result: res
                }
            });
        })
        .catch(e => {
            console.warn(e);
        });
}

export interface GetIndexedList extends Action {
    type: 'GET_INDEXED_LIST';
    payload: {
        result: Indexed[];
    };
}

export function getIndexedList(dispatch: Dispatch<RootState>): void {
    WebApi.get('indexed')
        .then((res: { result: Indexed[] }) => {
            // console.log(res);
            dispatch({
                type: 'GET_INDEXED_LIST',
                payload: res
            });
        })
        .catch(e => {
            console.warn(e);
        });
}

function _triggerSearch(filterParams?: FilterParams, page: number = 0): void {
    const queryParams = {
        ...filterParams,
        i: page
    };

    browserHistory.push(`/search?${WebApi.queryString(queryParams)}`);
}

export function undo() {
    return ActionCreators.undo();
}

export function redo() {
    return ActionCreators.redo();
}