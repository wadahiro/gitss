import * as React from 'react';

import { Tag } from './Tag';
import { Select, Option } from './Select';
import { Panel, PanelHeading, PanelBlock } from './Panel';
import { Menu, MenuLabel, MenuList, MenuLink } from './Menu';
import { Facet, FilterParams } from '../reducers';

interface FilterPanelProps {
    title: string;
    selected: string[];
    facet: Facet;
    onToggle: (term: string[]) => void;
    emptyKeyword?: string;
    emptyLabel?: string;
}


export class FilterPanel extends React.PureComponent<FilterPanelProps, void> {
    handleSelect = (values: Option[]) => {
        const terms = values.map(x => x.value);
        this.props.onToggle(terms);
    };

    render() {
        const { facet, title, emptyKeyword, emptyLabel, selected, onToggle } = this.props;

        if (!facet || !facet.terms || facet.terms.length === 0) {
            return null;
        }

        const style = {
            padding: 0
        };
        const options = facet.terms.map(x => {
            let type = x.term;
            if (emptyKeyword && emptyLabel && type === emptyKeyword) {
                type = emptyLabel;
            }
            return {
                label: `${type} (${x.count})`,
                value: x.term
            };
        });

        return (
            <Panel>
                <PanelHeading>{title}</PanelHeading>
                <Menu>
                    <MenuList>
                        <PanelBlock>
                            <li >
                                <Select options={options} multi value={selected} onChange={this.handleSelect} />
                            </li>
                        </PanelBlock>
                    </MenuList>
                </Menu>
            </Panel>
        );
    }
}

function contains(array = [], item) {
    return array !== null && typeof array.find(x => x === item) !== 'undefined';
}