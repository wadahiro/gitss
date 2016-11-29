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

interface PagerProps {
    pageSize: number;
    current: number;
    onChange: (page: number) => void;
}

export class Pager extends React.PureComponent<PagerProps, void> {
    showPage = (page: number) => {
        this.props.onChange(page);
    };

    next = () => {
        let next = this.props.current + 1;
        if (next >= this.props.pageSize) {
            next = this.props.pageSize - 1;
        }
        this.props.onChange(next);
    };

    prev = () => {
        let next = this.props.current - 1;
        if (next < 0) {
            next = 0;
        }
        this.props.onChange(next);
    };

    render() {
        const { pageSize, current } = this.props;

        const pageButtons = [];

        let start = current - 2;
        let end = current + 2;
        const lastIndex = pageSize - 1;

        if (start < 0) {
            end += -start;
            start = 0;
        }
        if (end > lastIndex) {
            if (start > 0) {
                start -= (end - lastIndex);
                if (start < 0) {
                    start = 0;
                }
            }
            end = lastIndex;
        }

        if (start > 0) {
            pageButtons.push((
                <li key={0}>
                    <PageButton onClick={this.showPage.bind(null, 0)}>1</PageButton>
                </li>
            ));
            pageButtons.push(<li key='dot-before'>...</li>);
        }

        for (let index = start; index <= end; index++) {
            pageButtons.push((
                <li key={index}>
                    <PageButton onClick={this.showPage.bind(null, index)} isActive={index === current}>{index + 1}</PageButton>
                </li>
            ));
        }

        if (end < lastIndex) {
            pageButtons.push(<li key='dot-after'>...</li>);
            pageButtons.push((
                <li key={lastIndex}>
                    <PageButton onClick={this.showPage.bind(null, lastIndex)}>{lastIndex + 1}</PageButton>
                </li>
            ));
        }

        return (
            <Pagination>
                <PageButton onClick={this.prev}>Previous</PageButton>
                <PageButton onClick={this.next}>Next page</PageButton>
                <ul>
                    {pageButtons}
                </ul>
            </Pagination>
        );
    }
}