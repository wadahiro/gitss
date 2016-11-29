import * as React from 'react';
import { Dispatch, Action } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router'

import { AppFooter } from '../components/Footer';
import { Grid, Container, Section, Row, Col } from '../components/Grid';
import { Hero, HeroHead, HeroBody, HeroFoot } from '../components/Hero';
import { Title, SubTitle } from '../components/Title';
import { InputTextAddon } from '../components/Input';

import { RootState, SearchResult, SearchFacets, BaseFilterParams, BaseFilterOptions, FilterParams, FacetKey } from '../reducers';
import * as Actions from '../actions';

interface Props {
    query: string;
    loading: boolean;
    showOptions: boolean;
    filterParams: FilterParams;
    result: SearchResult;
    facets: SearchFacets;
    // react-redux inject props
    dispatch?: Dispatch<Action>;
    // react-router inject props
    location?: any;
    history?: any;
    params?: BaseFilterParams;
    // lazy fetch
    baseFilterParams?: BaseFilterParams;
    baseFilterOptions?: BaseFilterOptions;
}

class HomeView extends React.PureComponent<Props, void> {

    handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.keyCode === 13) {
            Actions.triggerSearch(this.props.dispatch, this.props.baseFilterParams, e.target['value']);
        }
    };

    render() {
        const { loading, showOptions, query, filterParams, result, facets,
            baseFilterParams, baseFilterOptions } = this.props;

        return (
            <div style={viewStyle}>
                <div style={sectionStyle}>
                    <Section>
                    </Section>
                    <Hero color='isPrimary'>
                        <HeroBody>
                            <Container>
                                <Row>
                                    <Col size='is3'>
                                        <Title>GitSS</Title>
                                        <SubTitle>Git Source Search</SubTitle>
                                    </Col>
                                    <Col size='is8'>

                                        <InputTextAddon
                                            placeholder='Search'
                                            buttonIcon='fa fa-search'
                                            isExpanded={true}
                                            icon='fa fa-search'
                                            size='isLarge'
                                            hasIcon
                                            defaultValue={this.props.query}
                                            onKeyDown={this.handleKeyDown}
                                            />
                                    </Col>
                                </Row>
                            </Container>
                        </HeroBody>
                    </Hero>
                </div>
                <AppFooter />
            </div>
        );
    }
}

const viewStyle = {
    minHeight: '100vh',
    display: 'flex',
    flexDirection: 'column'
};

const sectionStyle = {
    flex: 1
};

const sideBarStyle = {
    flex: '0 0 300px',
    padding: 0
};

function mapStateToProps(state: RootState, props: Props): Props {
    return {
        loading: state.app.present.loading,
        showOptions: state.app.present.showSearchOptions,
        query: props.location.query['q'] !== undefined ? props.location.query['q'] : '',
        filterParams: props.location.query,
        result: state.app.present.result,
        facets: state.app.present.facets,

        // Convert react-router injected params
        baseFilterParams: toBaseFilterParams(props.params),
        baseFilterOptions: state.app.present.baseFilterOptions
    };
}

function toBaseFilterParams(params: Object): BaseFilterParams {
    return params;
}

const HomeViewContainer = connect(
    mapStateToProps
)(HomeView);

export default HomeViewContainer;