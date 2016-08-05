var gulp = require('gulp');
var sass = require('gulp-sass');
var t=require('./node_modules/t/index.js');


gulp.task('sass', function() {
    gulp.src('./input/sass/*.scss')
        .pipe(sass({
            sourceComments: true,
            outputStyle: 'expanded',
            errLogToConsole: true
        }))
        .pipe(gulp.dest('./output/css/'));
});
gulp.task('ttt', function() {
    gulp.src('./input/*.html')
        .pipe(t('./input/'))
        .pipe(gulp.dest('./output/'));
});




gulp.task('watch', ['sass','ttt'], function() {
    var sassWatcher = gulp.watch('./input/sass/*.scss', ['sass']);
    var tttWatcher = gulp.watch('./input/*.html', ['ttt']);
   
});