import * as React from 'react';

import { Size } from './Modifiers';
import { Title } from './Title';

const BLevel = require('re-bulma/lib/components/level/level').default;
const BLevelItem = require('re-bulma/lib/components/level/level-item').default;
const BHeading = require('re-bulma/lib/components/heading/heading').default;

interface LevelProps extends React.DOMAttributes {
}

export class Level extends React.PureComponent<LevelProps, void> {
    render() {
        return <BLevel {...this.props} />;
    }
}

interface LevelItemProps extends React.DOMAttributes {
    hasTextCentered?: boolean;
}

export class LevelItem extends React.PureComponent<LevelItemProps, void> {
    render() {
        return <BLevelItem {...this.props} />;
    }
}

export class Heading extends React.PureComponent<LevelProps, void> {
    render() {
        return <BHeading {...this.props} />;
    }
}

interface HeadingNavProps extends React.DOMAttributes {
    items: HeadingItem[];
}

interface HeadingItem {
    heading: string;
    title: string;
}

export class HeadingNav extends React.PureComponent<HeadingNavProps, void> {
    render() {
        return (
            <Level>
                {this.props.items.map(x => {
                    return (
                        <LevelItem key={x.heading} hasTextCentered>
                            <Heading>{x.heading}</Heading>
                            <Title>{x.title}</Title>
                        </LevelItem>
                    );
                })}
            </Level>
        );
    }
}