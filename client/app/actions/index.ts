import { Action, Dispatch } from 'redux';
import { browserHistory } from 'react-router';

const { ActionCreators } = require('redux-undo');

import { FilterParams } from '../reducers';
import WebApi from '../api/WebApi';

export type Actions =
    Search |
    SearchFilter |
    SearchStart
    ;

export interface Search extends Action {
    type: 'SEARCH';
    payload: {
        result: any;
    }
}

export interface SearchFilter extends Action {
    type: 'SEARCH_FILTER';
    payload: {
        result: any;
    }
}

export interface SearchStart extends Action {
    type: 'SEARCH_START';
    payload: {
        searchParams: FilterParams
    };
}

export function triggerSearch(dispatch: Dispatch<Search>, query?: string): void {
    // reset filters
    const params = {
        q: query
    };
    _triggerSearch(params, 0);
}

export function triggerFilter(dispatch: Dispatch<Search>, filterParams?: FilterParams, page: number = 0): void {
    _triggerSearch(filterParams, page);
}

export function search(dispatch: Dispatch<Search>, searchParams: FilterParams): void {

    dispatch({
        type: 'SEARCH_START',
        payload: {
            searchParams
        }
    });

    WebApi.query('search', searchParams)
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

function _triggerSearch(searchParams?: FilterParams, page: number = 0): void {
    const queryParams = {
        ...searchParams,
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