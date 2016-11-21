import { Action, Dispatch } from 'redux';
import { browserHistory } from 'react-router';

const { ActionCreators } = require('redux-undo');

import { BaseFilterParams, FilterParams } from '../reducers';
import WebApi from '../api/WebApi';

export type Actions =
    GetBaseFilters |
    Search |
    ResetFacets |
    SearchStart
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

export function triggerFilter(dispatch: Dispatch<Search>, filterParams?: FilterParams, page: number = 0): void {
    _triggerSearch(filterParams, page);
}

export function search(dispatch: Dispatch<Search>, baseFilterParams: BaseFilterParams = {}, filterParams: FilterParams): void {

    dispatch({
        type: 'SEARCH_START',
        payload: {
            filterParams
        }
    });

    const params = Object.keys(baseFilterParams).reduce((s, x) => {
        s[x[0]] = [baseFilterParams[x]];
        return s;
    }, filterParams);

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

function _triggerSearch(filterParams?: FilterParams, page: number = 0): void {
    const queryParams = {
        ...filterParams,
        i: page
    };

    browserHistory.push(`/?${WebApi.queryString(queryParams)}`);
}

export function undo() {
    return ActionCreators.undo();
}

export function redo() {
    return ActionCreators.redo();
}