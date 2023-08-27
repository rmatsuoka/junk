#!/usr/bin/env perl

use v5.36;
use experimental 'defer';
use Data::Dumper;

sub statdir($dir) {
	use File::stat;
	use File::Spec::Functions;

	my $prefix = $dir // "";
	$dir = $dir // ".";
	opendir my $dh, $dir or return;
	defer { closedir $dh; }

	my @stats;
	foreach (readdir $dh) {
		next if /^\.\.?$/;
		my $name = $prefix ? catfile($prefix, $_) : $_;
		my $stat = stat $name;
		push @stats, { stat => $stat, name => $name };
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
		my $stat = $_->{stat};
		my $rwx = format_mode $stat->mode;
		my $t = localtime $stat->mtime;
		printf(
			"%s %2d %s %s %10d %s %s\n",
			$rwx,
			$stat->nlink,
			scalar getgrgid $stat->gid,
			scalar getpwuid $stat->uid,
			$stat->size,
			$t->strftime("%b %d %Y %H:%M"),
			$_->{name}
		);
	}
}

unshift @ARGV, undef unless @ARGV;
foreach (@ARGV) {
	ls($_);
}
