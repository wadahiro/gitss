import * as React from 'react';

const BPagination = require('re-bulma/lib/components/pagination/pagination').default;
const BPageButton = require('re-bulma/lib/components/pagination/page-button').default;

interface PaginationProps extends React.HTMLAttributes {
    children?: React.ReactElement<any>;
    isActive?: boolean;
}

const defaultActiveStyle = {
    color: '#fff',
    backgroundColor: '#3572b0'
};

export class Pagination extends React.PureComponent<PaginationProps, void> {
    render() {
        return <BPagination {...this.props} >{this.props.children}</BPagination>;
    }
}

export class PageButton extends React.PureComponent<PaginationProps, void> {
    render() {
        let style = this.props.style;
        if (this.props.isActive) {
            style = {
                ...defaultActiveStyle,
                ...this.props.style
            };
        }

        return <BPageButton {...this.props} style={style}>{this.props.children}</BPageButton>;
    }
}
