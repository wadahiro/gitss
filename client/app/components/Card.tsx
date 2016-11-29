import * as React from 'react';

const BCard = require('re-bulma/lib/components/card/card').default;
const BCardHeader = require('re-bulma/lib/components/card/card-header').default;
const BCardHeaderTitle = require('re-bulma/lib/components/card/card-header-title').default;
const BCardHeaderIcon = require('re-bulma/lib/components/card/card-header-icon').default;
const BCardContent = require('re-bulma/lib/components/card/card-content').default;
const BCardFooter = require('re-bulma/lib/components/card/card-footer').default;
const BCardFooterItem = require('re-bulma/lib/components/card/card-footer-item').default;

interface CardPanelProps extends React.HTMLAttributes {
    title: string;
    icon: string;
    isFullWidth?: boolean;
    footers?: React.ReactElement<any>[];
    headerStyle?: any;
    onTitleToggle?: () => void;
}

const defaultStyle = {
};

export class CardPanel extends React.PureComponent<CardPanelProps, void> {
    static defaultProps = {
        headerStyle: {}
    };

    render() {
        const { onTitleToggle, ...rest } = this.props as any;
        return (
            <BCard style={defaultStyle} {...rest}>
                <BCardHeader style={this.props.headerStyle} onClick={this.props.onTitleToggle}>
                    <BCardHeaderTitle>
                        {this.props.title}
                    </BCardHeaderTitle>
                    <BCardHeaderIcon icon={this.props.icon} />
                </BCardHeader>
                <BCardContent>
                    {this.props.children}
                </BCardContent>
                {this.props.footers &&
                    <BCardFooter>
                        {this.props.footers.map(x => <BCardFooterItem>{x}</BCardFooterItem>)}
                    </BCardFooter>
                }
            </BCard>
        );
    }
}
