import * as React from 'react';
import { Dispatch, Action } from 'redux';
import { connect } from 'react-redux';

import { Tag } from './Tag';
import { FacetPanel } from '../components/FacetPanel';
import { RootState, SearchResult, SearchFacets, FilterParams, FilterParamKey } from '../reducers';
import * as Actions from '../actions';


interface SearchSidePanelProps {
    onToggle: (filterParams: FilterParams) => void;
    facets: SearchFacets;
    filterParams: FilterParams;
    query: string;
}

export class SearchSidePanel extends React.PureComponent<SearchSidePanelProps, void> {

    handleExtToggle = (term: string) => {
        this.props.onToggle(mergeTerm('x', this.props.filterParams, term));
    };

    handleOrganizationToggle = (term: string) => {
        this.props.onToggle(mergeTerm('o', this.props.filterParams, term));
    };

    handleProjectToggle = (term: string) => {
        this.props.onToggle(mergeTerm('p', this.props.filterParams, term));
    };

    handleRepositoryToggle = (term: string) => {
        this.props.onToggle(mergeTerm('r', this.props.filterParams, term));
    };

    handleBranchesToggle = (term: string) => {
        this.props.onToggle(mergeTerm('b', this.props.filterParams, term));
    };

    handleTagsToggle = (term: string) => {
        this.props.onToggle(mergeTerm('t', this.props.filterParams, term));
    };

    render() {
        const { facets,
            filterParams,
        } = this.props;

        return (
            <div>
                <FacetPanel title='File extensions'
                    facet={facets.facets['ext']}
                    emptyKeyword='/noext/'
                    emptyLabel='(No extension)'
                    selected={filterParams.x}
                    onToggle={this.handleExtToggle} />
                <FacetPanel title='Organization'
                    facet={facets.facets['organization']}
                    selected={filterParams.o}
                    onToggle={this.handleOrganizationToggle} />
                <FacetPanel title='Project'
                    facet={facets.facets['project']}
                    selected={filterParams.p}
                    onToggle={this.handleProjectToggle} />
                <FacetPanel title='Repository'
                    facet={facets.facets['repository']}
                    selected={filterParams.r}
                    onToggle={this.handleRepositoryToggle} />
                <FacetPanel title='Branches'
                    facet={facets.facets['branches']}
                    selected={filterParams.b}
                    onToggle={this.handleBranchesToggle} />
                <FacetPanel title='Tags'
                    facet={facets.facets['tags']}
                    selected={filterParams.t}
                    onToggle={this.handleTagsToggle} />

            </div>
        );
    }
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

