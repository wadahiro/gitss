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
    filterParams: FilterParams;
    query: string;
}

export class SearchSidePanel extends React.PureComponent<SearchSidePanelProps, void> {

    handleExtToggle = (terms: string[]) => {
        this.props.onToggle(mergeFilterParams('x', this.props.filterParams, terms));
    };

    handleOrganizationToggle = (terms: string[]) => {
        this.props.onToggle(mergeFilterParams('o', this.props.filterParams, terms));
    };

    handleProjectToggle = (terms: string[]) => {
        this.props.onToggle(mergeFilterParams('p', this.props.filterParams, terms));
    };

    handleRepositoryToggle = (terms: string[]) => {
        this.props.onToggle(mergeFilterParams('r', this.props.filterParams, terms));
    };

    handleBranchesToggle = (terms: string[]) => {
        this.props.onToggle(mergeFilterParams('b', this.props.filterParams, terms));
    };

    handleTagsToggle = (terms: string[]) => {
        this.props.onToggle(mergeFilterParams('t', this.props.filterParams, terms));
    };

    render() {
        const { facets,
            filterParams,
        } = this.props;

        let Panel = FilterPanel; // FacetPanel

        return (
            <div>
                <Panel title='File extensions'
                    facet={facets.facets['ext']}
                    emptyKeyword='/noext/'
                    emptyLabel='(No extension)'
                    selected={filterParams.x}
                    onToggle={this.handleExtToggle} />
                <Panel title='Organization'
                    facet={facets.facets['organization']}
                    selected={filterParams.o}
                    onToggle={this.handleOrganizationToggle} />
                <Panel title='Project'
                    facet={facets.facets['project']}
                    selected={filterParams.p}
                    onToggle={this.handleProjectToggle} />
                <Panel title='Repository'
                    facet={facets.facets['repository']}
                    selected={filterParams.r}
                    onToggle={this.handleRepositoryToggle} />
                <Panel title='Branches'
                    facet={facets.facets['branches']}
                    selected={filterParams.b}
                    onToggle={this.handleBranchesToggle} />
                <Panel title='Tags'
                    facet={facets.facets['tags']}
                    selected={filterParams.t}
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
    let terms;
    if (target.find(x => x === term)) {
        terms = target.filter(x => x !== term);
    } else {
        terms = target.concat(term);
    }
    return Object.assign({}, params, {
        [key]: terms
    });
}

