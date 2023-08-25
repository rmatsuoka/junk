#!/usr/bin/env perl

use v5.36;
use experimental 'defer';
use Data::Dumper;

sub statdir($dir) {
	use File::stat;
	use File::Spec::Functions;

	opendir my $dh, $dir or return;
	defer { closedir $dh; }

	my @stats;
	foreach (readdir $dh) {
		next if /^\.\.?$/;
		push @stats, { stat => stat $_, name => $_ };
	}
	return sort { $a->{name} cmp $b->{name} } @stats;
}

sub format_mode($perm) {
	my $rwx = 'rwxrwxrwx';
	my $l = length $rwx;
	for (0 .. $l - 1) {
		if (($perm & (1 << $_)) == 0) {
			substr $rwx, $l - $_ - 1, 1, '-';
		}
	}
	return $rwx
}

sub ls($dir) {
	use Time::Piece;

	foreach (statdir $dir) {
		local $, = ' ';
		my $stat = $_->{stat};
		my $rwx = format_mode $_->{stat}->mode;
		my $t = localtime $_->{stat}->mtime;
		printf(
			"%s %2d %s %s %10d %s %s\n",
			$rwx,
			$_->{stat}->nlink,
			scalar getgrgid $stat->gid,
			scalar getpwuid $stat->uid,
			$_->{stat}->size,
			$t->strftime("%b %d %Y %H:%M"),
			$_->{name}
		);
	}
}

foreach (@ARGV) {
	ls($_);
}
