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
    query: string;
    filterParams: FilterParams;
    lastQuery: string;
    facets: SearchFacets;
    result: SearchResult;
}

export interface FilterParams {
    ext: string[];
}

export interface SearchFacets {
    facets: Facets;
    fullRefsFacet: OranizationFacet[];
}

export interface SearchResult {
    query: string;
    filterParams: FilterParams;
    time: number;
    size: number;
    limit: number;
    current: number;
    next: number;
    isLastPage: boolean;
    hits: Hit[];
    facets?: Facets;
    fullRefsFacet?: OranizationFacet[];
}
export interface Hit {
    _source: Source;
    keyword: string[];
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

export interface OranizationFacet {
    term: string;
    count: number;
    projects: ProjectFacet[];
}

export interface ProjectFacet {
    term: string;
    count: number;
    repositories: RepositoryFace[];
}

export interface RepositoryFace {
    term: string;
    count: number;
    refs: RefFacets[];
}

export interface RefFacets {
    term: string;
    count: number;
}


function init(): AppState {
    return {
        loading: false,
        query: '',
        filterParams: {
            ext: []
        },
        lastQuery: '',
        facets: {
            facets: {},
            fullRefsFacet: []
        },
        result: {
            query: '',
            filterParams: {
                ext: []
            },
            time: -1,
            size: 0,
            limit: 0,
            current: 0,
            next: 0,
            isLastPage: true,
            hits: []
        }
    };
}

export const appStateReducer = (state: AppState = init(), action: Actions.Actions) => {
    switch (action.type) {
        case 'SET_QUERY':
            return Object.assign({}, state, {
                query: action.payload.query
            });
        case 'SEARCH_START':
            return Object.assign({}, state, {
                loading: true
            });
        case 'SEARCH':
            const searchResult: SearchResult = action.payload.result;

            let facets = {
                facets: searchResult.facets,
                fullRefsFacet: searchResult.fullRefsFacet
            };
            let filterParams = searchResult.filterParams;

            if (searchResult.query !== '' && searchResult.query === state.lastQuery) {
                // same query, so don't change facet view!
                facets = state.facets;
            } else {
                // search with new keyword
                filterParams = {
                    ext: []
                };
            }

            return Object.assign({}, state, {
                lastQuery: searchResult.query,
                filterParams,
                result: searchResult,
                facets,
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

