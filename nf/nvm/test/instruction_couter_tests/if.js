// Copyright (C) 2017 go-nebulas authors
//
// This file is part of the go-nebulas library.
//
// the go-nebulas library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-nebulas library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-nebulas library.  If not, see <http://www.gnu.org/licenses/>.
//
const assert = require('assert.js');

function assertEqual(func, args, expected, expected_count, msg) {
    const count_of_helper_statement = 9;
    var count = _instruction_counter.count;
    assert.equal(func.apply(null, args), expected);
    assert.equal(_instruction_counter.count - count - count_of_helper_statement, expected_count, msg);
};

var count = 0;
var x, y, z;

// test 1.
var test1 = function (x) {
    if (x < 5) {
        return 0;
    } else if (x < 20) {
        return 5;
    } else if (x < 40) {
        return 10;
    } else {
        return -1;
    }
};
assertEqual(test1, [1], 0, 1);
assertEqual(test1, [15], 5, 2);
assertEqual(test1, [30], 10, 3);
assertEqual(test1, [100], -1, 4);

// test 2.
var test2 = function (x, y) {
    if (x < 5 && x < y) {
        return 0;
    } else if (x < 50 && x < y) {
        return 50;
    } else if (x < y) {
        return 100;
    } else if (x > y) {
        return 200;
    } else {
        return 300;
    }
};
assertEqual(test2, [1, 3], 0, 3);
assertEqual(test2, [10, 15], 50, 6);
assertEqual(test2, [60, 70], 100, 7);
assertEqual(test2, [90, 80], 200, 8);
assertEqual(test2, [100, 100], 300, 8);

// test 3.
var test3 = function (x, y) {
    if (x < 5 || x < y) {
        return 0;
    } else if (x < 50 || x < y) {
        return 50;
    } else if (x < y) {
        return 100;
    } else if (x > y) {
        return 200;
    } else {
        return 300;
    }
};
assertEqual(test3, [1, 0], 0, 3);
assertEqual(test3, [5, 6], 0, 3);
assertEqual(test3, [10, 11], 0, 3);
assertEqual(test3, [40, 30], 50, 6);
assertEqual(test3, [60, 30], 200, 8);
assertEqual(test3, [60, 60], 300, 8);
