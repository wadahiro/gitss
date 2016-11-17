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
    facets: SearchFacets;
    result: SearchResult;
}

export interface FilterParams {
    q?: string;
    i?: number;
    a?: AdvancedSearchType;
    x?: string[]; // ext
    o?: string[]; // organization
    p?: string[]; // project
    r?: string[]; // repository
    b?: string[]; // branches
    t?: string[]; // tags
}

export type AdvancedSearchType = 'regex';

export type FilterParamKey = 'q' | 'x' | 'o' | 'p' | 'r' | 'b' | 't';

const FILTER_PARAMS_MAP: { [index: string]: FacetKey } = {
    x: 'ext',
    o: 'organization',
    p: 'project',
    r: 'repository',
    b: 'branches',
    t: 'tags'
};

export interface SearchFacets {
    facets: Facets;
    fullRefsFacet: OranizationFacet[];
}

export interface SearchResult {
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

export interface Hit extends FileMetadata {
    keyword: string[];
    preview: Preview[];
}

export interface Preview {
    offset: number;
    preview: string;
    hits: number[];
}

export interface FileMetadata {
    blob: string;
    organization: string;
    project: string;
    repository: string;
    branches: string[];
    tags: string[];
    path: string;
    ext: string;
}

export type FacetKey = 'ext' | 'organization' | 'project' | 'repository' | 'branches' | 'tags';

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
        facets: {
            facets: {},
            fullRefsFacet: []
        },
        result: {
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
        case 'SEARCH_START':
            return {
                ...state,
                loading: true,
                searchParams: action.payload.searchParams || {}
            }
        case 'SEARCH':
        case 'SEARCH_FILTER':
            const searchResult: SearchResult = action.payload.result;

            let facets = {
                facets: searchResult.facets,
                initialFacets: {},
                fullRefsFacet: searchResult.fullRefsFacet
            };
            // let filterParams = searchResult.filterParams;

            // if (action.type === 'SEARCH_FILTER') {
            //     // same query, so don't reduce facet items. we need to update the values.

            //     facets.facets = Object.keys(FILTER_PARAMS_MAP).reduce((s, k) => {
            //         const noSeleted = filterParams[k] === undefined;
            //         const facetKey = FILTER_PARAMS_MAP[k];
            //         s[facetKey] = {
            //             ...state.facets.facets[facetKey],
            //             terms: mergeTerms(state.facets.facets, searchResult.facets, facetKey, noSeleted)
            //         };
            //         return s;
            //     }, {} as Facets);
            // } else {
            //     // search with new keyword
            //     filterParams = {};
            // }

            window.scrollTo(0, 0);

            return {
                ...state,
                result: searchResult,
                facets,
                loading: false
            }
    }

    return state;
};

function mergeTerms(prev: Facets, next: Facets, key: string, noSeleted: boolean): Term[] {
    const prevTerms = prev[key] ? prev[key].terms : [];
    const nextTerms = next[key] ? next[key].terms : [];
    return prevTerms.map(x => {
        const nextTerm = nextTerms.find(y => x.term === y.term);
        if (nextTerm) {
            if (noSeleted) {
                x.count = nextTerm.count;
            }
        } else {
            if (noSeleted) {
                x.count = 0;
            }
        }
        return x;
    });
}

export default combineReducers({
    app: undoable(appStateReducer, {
        filter: includeAction(['ADD_ITEM', 'DELETE_ITEM', 'MOD_QUANTITY', 'MOD_EXCHANGE_RATE', 'MOD_METADATA', 'RESTORE_SAVED_HISTORY'])
    }),
});

