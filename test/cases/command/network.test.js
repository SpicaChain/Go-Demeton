'use strict';

var process = require('child_process');
var expect = require('chai').expect;
var FS = require('fs');
var deb = './deb';
var newAccount;

function countSubstr(str, substr) {
    var reg = new RegExp(substr, "g");
    return str.match(reg) ? str.match(reg).length : 0;//若match返回不为null，则结果为true，输出match返回的数组(["test","test"])的长度  
}

describe('deb network', () => {

    it('deb network ssh-keygen', (done) => {
        process.exec(deb + ' network ssh-keygen ./network.tmp', (err, stdout, stderr) => {
            try {
                expect(err).to.be.equal(null);
                expect(stderr).to.be.equal("");

                var network = FS.readFileSync('./network.tmp','utf-8');
                var expectNetwork = /.+==/
                
                expect(expectNetwork.test(network)).to.be.equal(true);

            } catch (err) {
                done(err);
                return;
            }
            done();
        });
    });

    after('rm network.tmp', (done) => {
        process.exec('rm -f ./network.tmp', (err, stdout, stderr) => {
            try {
                expect(err).to.be.equal(null);
                expect(stderr).to.be.equal("");
            } catch (err) {
                done(err);
                return;
            }
            done();
        })
    })
});
