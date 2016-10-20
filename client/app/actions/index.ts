import { Action, Dispatch } from 'redux';

const { ActionCreators } = require('redux-undo');

import WebApi from '../api/WebApi';

export type Actions =
    Search |
    SearchStart
    ;

export interface Search extends Action {
    type: 'SEARCH';
    payload: {
        result: any;
    }
}

export interface SearchStart extends Action {
    type: 'SEARCH_START';
}

export function search(dispatch: Dispatch<Search>, query: string): void {
    dispatch({
        type: 'SEARCH_START'
    });

    WebApi.get('search?q=' + query)
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