import { combineReducers } from 'redux';
import { Maybe, Either } from 'tsmonad';
import * as lm from 'lodash/mergeWith';
import * as lu from 'lodash/unionWith';

import * as Actions from '../actions';

const ReduxUndo = require('redux-undo');
const undoable = ReduxUndo.default;
const includeAction = ReduxUndo.includeAction;
const mergeWith = lm['default'];
const unionWith = lu['default'];


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
    showSearchOptions: boolean;

    facets: SearchFacets;
    result: SearchResult;

    indexedList: Indexed[]

    baseFilterOptions: BaseFilterOptions;
}

export interface Indexed extends RepositoryMetadata {
    lastUpdated: string; // YYYY-MM-DD HH:mm:dd Z
    branches: {
        [index: string]: string;
    };
    tags: {
        [index: string]: string;
    };
}

export interface BaseFilterParams {
    organization?: string;
    project?: string;
    repository?: string;
    branch?: string;
    tag?: string;
}

export interface BaseFilterOptions {
    organizations: Option[];
    projects: Option[];
    repositories: Option[];
    branches: Option[];
    tags: Option[];
}

export interface Option {
    label: string;
    value: string;
}

export interface FilterParams {
    a?: AdvancedSearchType;
    x?: string[]; // ext
    o?: string[]; // organization
    p?: string[]; // project
    r?: string[]; // repository
    b?: string[]; // branches
    t?: string[]; // tags
}

export type AdvancedSearchType = 'regex';

export type FilterParamKey = 'x' | 'o' | 'p' | 'r' | 'b' | 't';

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

export interface FileMetadata extends RepositoryMetadata {
    blob: string;
    organization: string;
    project: string;
    repository: string;
    branches: string[];
    tags: string[];
    path: string;
    ext: string;
}

export interface RepositoryMetadata {
    organization: string;
    project: string;
    repository: string;
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
    repositories: RepositoryFacet[];
}

export interface RepositoryFacet {
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
        showSearchOptions: false,
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
        },
        baseFilterOptions: {
            organizations: [],
            projects: [],
            repositories: [],
            branches: [],
            tags: []
        },
        indexedList: []
    };
}

function toOptions(array: string[] = []): Option[] {
    return array.map(x => ({ label: x, value: x }));
}

export const appStateReducer = (state: AppState = init(), action: Actions.Actions) => {
    switch (action.type) {
        case 'TOGGLE_SEARCH_OPTIONS':
            return {
                ...state,
                showSearchOptions: !state.showSearchOptions
            };

        case 'GET_INDEXED_LIST':
            return {
                ...state,
                indexedList: action.payload.result
            };

        case 'GET_BASE_FILTERS':

            return {
                ...state,
                baseFilterOptions: {
                    organizations: toOptions(action.payload.organizations),
                    projects: toOptions(action.payload.projects),
                    repositories: toOptions(action.payload.repositories),
                    branches: toOptions(action.payload.branches),
                    tags: toOptions(action.payload.tags)
                } as BaseFilterOptions
            };

        case 'SEARCH_START':
            return {
                ...state,
                loading: true,
                filterParams: action.payload.filterParams || {}
            };

        case 'RESET_FACETS':
            return {
                ...state,
                facets: {
                    facets: {},
                    fullRefsFacet: []
                }
            };

        case 'SEARCH':
            const searchResult: SearchResult = action.payload.result;

            const merged = mergeWith(state.facets.facets, searchResult.facets, (objValue: Facet, srcValue: Facet) => {
                const mergedFacet = mergeWith(objValue, srcValue, (objValue: any, srcValue: any) => {
                    if (Array.isArray(objValue)) {
                        return unionWith(objValue, srcValue, (a: Term, b: Term) => {
                            return a.term === b.term;
                        });
                    }
                    return srcValue;
                })
                return mergedFacet;
            });

            let facets = {
                facets: merged,
                fullRefsFacet: searchResult.fullRefsFacet
            };

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

export default combineReducers({
    app: undoable(appStateReducer, {
        filter: includeAction(['ADD_ITEM', 'DELETE_ITEM', 'MOD_QUANTITY', 'MOD_EXCHANGE_RATE', 'MOD_METADATA', 'RESTORE_SAVED_HISTORY'])
    }),
});

