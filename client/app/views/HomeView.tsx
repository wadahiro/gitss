import * as React from 'react';
import { Dispatch, Action } from 'redux';
import { connect } from 'react-redux';
import { Link, RouteComponentProps } from 'react-router';

import { StickyContainer, Sticky } from '../components/Sticky';
import { AppFooter } from '../components/Footer';
import { Grid, Container, Section, Row, Col, StickyFooterPage } from '../components/Grid';
import { Hero, HeroHead, HeroBody, HeroFoot } from '../components/Hero';
import { Title, SubTitle } from '../components/Title';
import { InputTextAddon } from '../components/Input';
import { MediaPanel } from '../components/Media';
import { HeadingNav } from '../components/Level';
import { RefTag } from '../components/RefTag';

import { RootState, SearchResult, SearchFacets, BaseFilterParams, BaseFilterOptions, FilterParams, FacetKey, Indexed } from '../reducers';
import * as Actions from '../actions';

interface Params {
}

interface Props extends RouteComponentProps<Params, void> {
    indexedList: Indexed[];
    // react-redux inject props
    dispatch?: Dispatch<Action>;
}

class HomeView extends React.PureComponent<Props, void> {
    handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.keyCode === 13) {
            Actions.triggerSearch(this.props.dispatch, e.target['value']);
        }
    };

    handleClick = (value) => {
        Actions.triggerSearch(this.props.dispatch, value);
    };

    componentWillMount() {
        Actions.getIndexedList(this.props.dispatch);
    }

    render() {
        return (
            <StickyContainer>
                <StickyFooterPage footer={<AppFooter />}>
                    <Section>
                        <IndexedSummary indexedList={this.props.indexedList} />
                    </Section>
                    <Sticky stickyStyle={{ zIndex: 1 }}>
                        <Hero color='isPrimary'>
                            <HeroBody>
                                <Container>
                                    <Row>
                                        <Col size='is3'>
                                            <Title>GitSS</Title>
                                            <SubTitle>Git Source Search</SubTitle>
                                        </Col>
                                        <Col size='is6'>
                                            <InputTextAddon
                                                placeholder='Search'
                                                buttonIcon='fa fa-search'
                                                onButtonClick={this.handleClick}
                                                isExpanded={true}
                                                icon='fa fa-search'
                                                size='isLarge'
                                                hasIcon
                                                onKeyDown={this.handleKeyDown}
                                                />
                                        </Col>
                                    </Row>
                                </Container>
                            </HeroBody>
                        </Hero>
                    </Sticky>
                    <Section>
                        <Container>
                            {this.props.indexedList.map(x => {
                                return (
                                    <MediaPanel key={x.lastUpdated} icon='feed'>
                                        <p>
                                            <h4>{x.lastUpdated}</h4>
                                            <h5><strong>{x.organization} : {x.project}/{x.repository}</strong> was synced.</h5>

                                            <h6>Branches:</h6>
                                            <p>
                                                {Object.keys(x.branches).map(k => {
                                                    return (
                                                        <RefTag key={k} type='branch'>{k}</RefTag>
                                                    );
                                                })}
                                            </p>

                                            <h6>Tags:</h6>
                                            <p>
                                                {Object.keys(x.tags).length === 0 ?
                                                    <span>Nothing.</span>
                                                    :
                                                    Object.keys(x.tags).map(k => {
                                                        return (
                                                            <RefTag key={k} type='tag'>{k}</RefTag>
                                                        );
                                                    })
                                                }
                                            </p>
                                        </p>
                                    </MediaPanel>
                                );
                            })}
                        </Container>
                    </Section>
                </StickyFooterPage>
            </StickyContainer>
        );
    }
}

interface IndexedSummaryProps {
    indexedList: Indexed[];
}

class IndexedSummary extends React.PureComponent<IndexedSummaryProps, any> {
    render() {
        const { indexedList } = this.props;
        const items = [
            {
                heading: 'Organizations',
                title: indexedList.length.toString()
            },
            {
                heading: 'Projects',
                title: indexedList.length.toString()
            },
            {
                heading: 'Repositories',
                title: indexedList.length.toString()
            },
            {
                heading: 'Branches',
                title: indexedList.length.toString()
            },
            {
                heading: 'Tags',
                title: indexedList.length.toString()
            },
            {
                heading: 'Files',
                title: indexedList.length.toString()
            }
        ];
        return (
            <HeadingNav items={items} />
        );
    }
}

function mapStateToProps(state: RootState, props: Props): Props {
    return {
        indexedList: state.app.present.indexedList
    };
}

const HomeViewContainer = connect(
    mapStateToProps
)(HomeView);

export default HomeViewContainer;