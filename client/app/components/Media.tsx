import * as React from 'react';

import { Size } from './Modifiers';
import { Icon } from './Icon';
import { Content } from './Grid';

const BMedia = require('re-bulma/lib/components/media/media').default;
const BMediaContent = require('re-bulma/lib/components/media/media-content').default;
const BMediaLeft = require('re-bulma/lib/components/media/media-left').default;
const BMediaRight = require('re-bulma/lib/components/media/media-right').default;


interface MediaProps extends React.HTMLAttributes {
    color?: 'isPrimary' | 'isDanger';
}

export class Media extends React.PureComponent<MediaProps, void> {
    render() {
        return <BMedia {...this.props} />;
    }
}

export class MediaContent extends React.PureComponent<MediaProps, void> {
    render() {
        return <BMediaContent {...this.props} />;
    }
}

export class MediaLeft extends React.PureComponent<MediaProps, void> {
    render() {
        return <BMediaLeft {...this.props} />;
    }
}

export class MediaRight extends React.PureComponent<MediaProps, void> {
    render() {
        return <BMediaRight {...this.props} />;
    }
}

interface MediaPanelProps extends React.HTMLAttributes {
    icon: string;
}

export class MediaPanel extends React.PureComponent<MediaPanelProps, void> {
    render() {
        return (
            <Media>
                <MediaLeft>
                    <Icon icon={this.props.icon} size='isMedium' />
                </MediaLeft>
                <MediaContent>
                    <Content>
                        {this.props.children}
                    </Content>
                </MediaContent>
            </Media>
        );
    }
}