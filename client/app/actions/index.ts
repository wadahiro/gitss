import { Action, Dispatch } from 'redux';

const { ActionCreators } = require('redux-undo');

import WebApi from '../api/WebApi';

export type Actions =
    Search
    ;

export interface Search extends Action {
    type: 'SEARCH';
    payload: {
        result: any;
    }
}
export function search(dispatch: Dispatch<Search>, query: string): void {
    WebApi.get('search?q=' + query)
        .then(res => {
            console.log(res);
            dispatch({
                type: 'SEARCH',
                payload: {
                    result: res
                }
            });
        })
        .catch(e => {
            console.warn(e);
        })
}

export function undo() {
    return ActionCreators.undo();
}

export function redo() {
    return ActionCreators.redo();
}