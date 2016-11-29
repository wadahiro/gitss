import * as React from 'react';
import { Maybe } from 'tsmonad';
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

import { RootState, SearchResult, SearchFacets, BaseFilterParams, BaseFilterOptions, FilterParams, FacetKey, Statistics } from '../reducers';
import * as Actions from '../actions';

interface Params {
}

interface Props extends RouteComponentProps<Params, void> {
    statistics: Maybe<Statistics>;
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
        Actions.getStatistics(this.props.dispatch);
    }

    render() {
        return (
            <StickyContainer>
                <StickyFooterPage footer={<AppFooter />}>
                    <Section>
                        {this.props.statistics.caseOf({
                            just: statistics => <StatisticsNav statistics={statistics} />,
                            nothing: () => <StatisticsNav loading />
                        })}
                    </Section>
                    <Sticky stickyStyle={{ zIndex: 1 }}>
                        <Hero color='isPrimary'>
                            <HeroBody style={{ padding: '20px 20px' }}>
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
                                        <Col size='is3'>
                                        </Col>
                                    </Row>
                                </Container>
                            </HeroBody>
                        </Hero>
                    </Sticky>
                    <Section>
                        <Container>
                            {this.props.statistics.caseOf({
                                just: statistics => {
                                    return statistics.indexes.map(x => {
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
                                    })
                                },
                                nothing: () => null
                            })}
                        </Container>
                    </Section>
                </StickyFooterPage>
            </StickyContainer>
        );
    }
}

interface StatiStatisticsProps {
    loading?: boolean;
    statistics?: Statistics;
}

class StatisticsNav extends React.PureComponent<StatiStatisticsProps, any> {
    format(key: string) {
        if (this.props.loading) {
            return '...';
        }
        const num = this.props.statistics.count[key];
        return String(num).replace(/(\d)(?=(\d\d\d)+(?!\d))/g, '$1,');
    }

    render() {
        const { statistics } = this.props;

        const items = [
            {
                heading: 'Organizations',
                title: this.format('organization')
            },
            {
                heading: 'Projects',
                title: this.format('project')
            },
            {
                heading: 'Repositories',
                title: this.format('repository')
            },
            {
                heading: 'Branches',
                title: this.format('branch')
            },
            {
                heading: 'Tags',
                title: this.format('tag')
            },
            {
                heading: 'Documents',
                title: this.format('document')
            }
        ];

        return (
            <HeadingNav items={items} />
        );
    }
}

function mapStateToProps(state: RootState, props: Props): Props {
    return {
        statistics: state.app.present.statistics
    };
}

const HomeViewContainer = connect(
    mapStateToProps
)(HomeView);

export default HomeViewContainer;