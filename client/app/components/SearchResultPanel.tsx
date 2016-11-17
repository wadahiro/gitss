import * as React from 'react';
import { Dispatch, Action } from 'redux';
import { connect } from 'react-redux';
import LazyLoad from 'react-lazyload';

import { Grid, Section, Row, Col } from '../components/Grid';
import { FileContent } from '../components/FileContent';
import { Pager } from '../components/Pagination';
import { RootState, SearchResult, SearchFacets, FilterParams, FacetKey } from '../reducers';

interface SearchResultPanelProps {
    result: SearchResult;
    onPageChange: (page: number) => void;
}

export class SearchResultPanel extends React.PureComponent<SearchResultPanelProps, void> {
    render() {
        const { result } = this.props;
        const pageSize = Math.ceil(result.size / result.limit) || 0;

        return (
            <Grid>
                {result && result.size > 10 &&
                    <Row>
                        <Col size='is12'>
                            <Pager pageSize={pageSize} current={result.current} onChange={this.props.onPageChange} />
                        </Col>
                    </Row>
                }
                {result.hits.map(x => {
                    return (
                        <Row key={x.blob}>
                            <LazyLoad height={200}>
                                <Col size='is12'>
                                    <FileContent metadata={x} keyword={x.keyword} preview={x.preview} />
                                </Col>
                            </LazyLoad>
                        </Row>
                    );
                })}
                {result && result.size > 10 &&
                    <Row>
                        <Col size='is12'>
                            <Pager pageSize={pageSize} current={result.current} onChange={this.props.onPageChange} />
                        </Col>
                    </Row>
                }
            </Grid>
        );
    }
}
