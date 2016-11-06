import { Action, Dispatch } from 'redux';

const { ActionCreators } = require('redux-undo');

import { FilterParams } from '../reducers';
import WebApi from '../api/WebApi';

export type Actions =
    SetQuery |
    Search |
    SearchFilter |
    SearchStart
    ;

export interface SetQuery extends Action {
    type: 'SET_QUERY';
    payload: {
        query: string;
    }
}

export function setQuery(dispatch: Dispatch<SetQuery>, query: string) {
    dispatch({
        type: 'SET_QUERY',
        payload: {
            query
        }
    });
}

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
        filterParams: FilterParams
    };
}

export function search(dispatch: Dispatch<Search>, query: string, filterParams?: FilterParams, page: number = 0): void {
    _search('SEARCH', dispatch, query, filterParams, page);
}

export function searchFilter(dispatch: Dispatch<Search>, query: string, filterParams?: FilterParams, page: number = 0): void {
    _search('SEARCH_FILTER', dispatch, query, filterParams, page);
}

function _search(searchType: 'SEARCH' | 'SEARCH_FILTER', dispatch: Dispatch<Search>, query: string, filterParams?: FilterParams, page: number = 0): void {
    dispatch({
        type: 'SEARCH_START',
        payload: {
            filterParams
        }
    });

    const queryParams = Object.assign({}, filterParams, {
        q: query,
        i: page
    });

    WebApi.query('search', queryParams)
        .then(res => {
            // console.log(res);
            dispatch({
                type: searchType,
                payload: {
                    result: res
                }
            });
        })
        .catch(e => {
            console.warn(e);
        });
}

export function undo() {
    return ActionCreators.undo();
}

export function redo() {
    return ActionCreators.redo();
}