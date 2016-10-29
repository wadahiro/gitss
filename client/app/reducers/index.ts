import { combineReducers } from 'redux';
import { Maybe, Either } from 'tsmonad';


import * as Actions from '../actions';

const ReduxUndo = require('redux-undo');
const undoable = ReduxUndo.default;
const includeAction = ReduxUndo.includeAction;


export interface RootState {
    app: AppStateHistory;
}

export interface AppStateHistory {
    past: AppState[];
    present: AppState;
    future: AppState[];
}

export interface AppState {
    loading: boolean;
    result: SearchResult;
}

export interface SearchResult {
    time: number;
    size: number;
    limit: number;
    current: number;
    next: number;
    isLastPage: boolean;
    hits: Hit[];
    facets: Facets;
}
export interface Hit {
    _source: Source;
    preview: Preview[];
}
export interface Preview {
    offset: number;
    preview: string;
    hits: number[];
}
export interface Source {
    blob: string;
    content: string;
    metadata: FileMetadata;
}

export interface FileMetadata {
    organization: string;
    project: string;
    repository: string;
    refs: string[];
    path: string;
    ext: string;
}

export interface Facets {
    [index: string]: Facet;
}

export interface Facet {
    field: string;
    missing: number;
    other: number;
    total: number;
    terms: Term[];
}

export interface Term {
    term: string;
    count: number;
}

function init(): AppState {
    return {
        loading: false,
        result: {
            time: -1,
            size: 0,
            limit: 0,
            current: 0,
            next: 0,
            isLastPage: true,
            hits: [],
            facets: {}
        }
    };
}

export const appStateReducer = (state: AppState = init(), action: Actions.Actions) => {
    switch (action.type) {
        case 'SEARCH_START':
            return Object.assign({}, state, {
                loading: true
            });
        case 'SEARCH':
            return Object.assign({}, state, {
                result: action.payload.result,
                loading: false
            });
    }

    return state;
};

export default combineReducers({
    app: undoable(appStateReducer, {
        filter: includeAction(['ADD_ITEM', 'DELETE_ITEM', 'MOD_QUANTITY', 'MOD_EXCHANGE_RATE', 'MOD_METADATA', 'RESTORE_SAVED_HISTORY'])
    }),
});
