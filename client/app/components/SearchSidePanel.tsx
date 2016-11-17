import * as React from 'react';
import { Dispatch, Action } from 'redux';
import { connect } from 'react-redux';

import { Tag } from './Tag';
import { FacetPanel } from './FacetPanel';
import { FilterPanel } from './FilterPanel';
import { RootState, SearchResult, SearchFacets, FilterParams, FilterParamKey } from '../reducers';
import * as Actions from '../actions';


interface SearchSidePanelProps {
    onToggle: (filterParams: FilterParams) => void;
    facets: SearchFacets;
    searchParams: FilterParams;
}

export class SearchSidePanel extends React.PureComponent<SearchSidePanelProps, void> {

    handleExtToggle = (terms: string[]) => {
        this.props.onToggle(mergeFilterParams('x', this.props.searchParams, terms));
    };

    handleOrganizationToggle = (terms: string[]) => {
        this.props.onToggle(mergeFilterParams('o', this.props.searchParams, terms));
    };

    handleProjectToggle = (terms: string[]) => {
        this.props.onToggle(mergeFilterParams('p', this.props.searchParams, terms));
    };

    handleRepositoryToggle = (terms: string[]) => {
        this.props.onToggle(mergeFilterParams('r', this.props.searchParams, terms));
    };

    handleBranchesToggle = (terms: string[]) => {
        this.props.onToggle(mergeFilterParams('b', this.props.searchParams, terms));
    };

    handleTagsToggle = (terms: string[]) => {
        this.props.onToggle(mergeFilterParams('t', this.props.searchParams, terms));
    };

    render() {
        const { facets,
            searchParams,
        } = this.props;

        if (!facets || !facets.facets) {
            return null;
        }

        let Panel = FilterPanel; // FacetPanel

        return (
            <div>
                <Panel title='File extensions'
                    facet={facets.facets['ext']}
                    emptyKeyword='/noext/'
                    emptyLabel='(No extension)'
                    selected={searchParams.x}
                    onToggle={this.handleExtToggle} />
                <Panel title='Organization'
                    facet={facets.facets['organization']}
                    selected={searchParams.o}
                    onToggle={this.handleOrganizationToggle} />
                <Panel title='Project'
                    facet={facets.facets['project']}
                    selected={searchParams.p}
                    onToggle={this.handleProjectToggle} />
                <Panel title='Repository'
                    facet={facets.facets['repository']}
                    selected={searchParams.r}
                    onToggle={this.handleRepositoryToggle} />
                <Panel title='Branches'
                    facet={facets.facets['branches']}
                    selected={searchParams.b}
                    onToggle={this.handleBranchesToggle} />
                <Panel title='Tags'
                    facet={facets.facets['tags']}
                    selected={searchParams.t}
                    onToggle={this.handleTagsToggle} />

            </div>
        );
    }
}

function mergeFilterParams(key: FilterParamKey, prev: FilterParams, terms: string[]): FilterParams {
    return {
        ...prev,
        [key]: terms
    };
}

function mergeTerm(key: FilterParamKey, params: FilterParams, term: string) {
    const target = params[key] || [];
    if (Array.isArray(target)) {
        let terms;
        if (target.find(x => x === term)) {
            terms = target.filter(x => x !== term);
        } else {
            terms = target.concat(term);
        }
        return {
            ...params,
            [key]: terms
        };
    } else {
        return {
            ...params,
            [key]: term
        };
    }
}

