import { Action, Dispatch } from 'redux';
import { browserHistory } from 'react-router';

const { ActionCreators } = require('redux-undo');

import { RootState, BaseFilterParams, FilterParams } from '../reducers';
import WebApi from '../api/WebApi';

export type Actions =
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

type BaseFiltersType = 'organization' | 'project' | 'repository' | 'branch' | 'tag';

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

export function triggerSearch(dispatch: Dispatch<Search>, baseFilterParmas: BaseFilterParams, query?: string): void {
    // reset filters & current facets
    dispatch({
        type: 'RESET_FACETS'
    });

    const params = {
        q: query
    };
    _triggerSearch(baseFilterParmas, params, 0);
}

export function triggerFilter(dispatch: Dispatch<Search>, baseFilterParmas: BaseFilterParams, filterParams: FilterParams, query: string, page: number = 0): void {
    const params = {
        ...filterParams,
        q: query
    };
    _triggerSearch(baseFilterParmas, params, page);
}

export function triggerBaseFilter(dispatch: Dispatch<Search>, baseFilterParams: BaseFilterParams, filterParams: FilterParams, query: string): void {
    browserHistory.push(`${_makeBaseFilterPath(baseFilterParams)}?${_makeQueryString(filterParams, query)}`);
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

function _makeBaseFilterPath(baseFilterParams: BaseFilterParams) {
    let urlPath = '/';
    if (baseFilterParams.organization) {
        urlPath += `s/${baseFilterParams.organization}`;
        if (baseFilterParams.project) {
            urlPath += `/${baseFilterParams.project}`;
            if (baseFilterParams.repository) {
                urlPath += `/${baseFilterParams.repository}`;
                if (baseFilterParams.branch) {
                    urlPath += `/branches/${baseFilterParams.branch}`;
                } else if (baseFilterParams.tag) {
                    urlPath += `/tags/${baseFilterParams.tag}`;
                }
            }
        }
    }
    return urlPath;
}

export function search(dispatch: Dispatch<Search>, baseFilterParams: BaseFilterParams = {}, filterParams: FilterParams): void {

    dispatch({
        type: 'SEARCH_START',
        payload: {
            filterParams
        }
    });

    const params = {
        ...filterParams,
        ..._toShortKeyParams(baseFilterParams)
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

function _triggerSearch(baseFilterParams: BaseFilterParams, filterParams?: FilterParams, page: number = 0): void {
    const queryParams = {
        ...filterParams,
        i: page
    };

    browserHistory.push(`/search?${WebApi.queryString(queryParams)}`);
}

function _toShortKeyParams(baseFilterParams: BaseFilterParams) {
    const params = Object.keys(baseFilterParams).reduce((s, x) => {
        if (typeof baseFilterParams[x] === 'object') {
            s[baseFilterParams[x].name[0]] = [baseFilterParams[x].value];
        } else {
            s[x[0]] = [baseFilterParams[x]];
        }
        return s;
    }, {});
    return params;
}


export function undo() {
    return ActionCreators.undo();
}

export function redo() {
    return ActionCreators.redo();
}