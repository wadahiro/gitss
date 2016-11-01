import { Action, Dispatch } from 'redux';

const { ActionCreators } = require('redux-undo');

import { FilterParams } from '../reducers';
import WebApi from '../api/WebApi';

export type Actions =
    SetQuery |
    Search |
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

export interface SearchStart extends Action {
    type: 'SEARCH_START';
}

export function search(dispatch: Dispatch<Search>, query: string, filterParams?: FilterParams): void {
    dispatch({
        type: 'SEARCH_START'
    });

    const queryParams = Object.assign({}, filterParams, {
        q: query
    });

    WebApi.query('search', queryParams)
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

export function undo() {
    return ActionCreators.undo();
}

export function redo() {
    return ActionCreators.redo();
}