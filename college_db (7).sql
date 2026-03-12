-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1:3307
-- Generation Time: Mar 11, 2026 at 06:37 AM
-- Server version: 10.4.32-MariaDB
-- PHP Version: 8.2.12

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `college_db`
--

-- --------------------------------------------------------

--
-- Table structure for table `admin_users`
--

CREATE TABLE `admin_users` (
  `id` int(11) NOT NULL,
  `register_no` varchar(20) NOT NULL,
  `password` varchar(255) NOT NULL,
  `department` varchar(50) NOT NULL,
  `name` varchar(250) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `admin_users`
--

INSERT INTO `admin_users` (`id`, `register_no`, `password`, `department`, `name`) VALUES
(1, '1254', 'c55b3f6a81a515bffa332d727560f1b5', 'CSE', 'Vignesh Saravanan'),
(0, '1255', 'c55b3f6a81a515bffa332d727560f1b5', 'CSE', 'Vijayalakshmi B'),
(2, '1257', 'c55b3f6a81a515bffa332d727560f1b5', 'CSE', 'Testing');

-- --------------------------------------------------------

--
-- Table structure for table `cas`
--

CREATE TABLE `cas` (
  `id` int(11) NOT NULL,
  `register_no` varchar(50) NOT NULL,
  `name` varchar(100) NOT NULL,
  `password` varchar(255) NOT NULL,
  `department` varchar(50) NOT NULL,
  `year` int(11) DEFAULT NULL,
  `section` varchar(10) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `cas`
--

INSERT INTO `cas` (`id`, `register_no`, `name`, `password`, `department`, `year`, `section`) VALUES
(1, '1204', 'Swarna Sudha S', 'c55b3f6a81a515bffa332d727560f1b5', 'CSE', 3, 'A'),
(2, '1345', 'Vijaya Amala Devi', 'c55b3f6a81a515bffa332d727560f1b5', 'CSE', 3, 'B');

-- --------------------------------------------------------

--
-- Table structure for table `hods`
--

CREATE TABLE `hods` (
  `id` int(11) NOT NULL,
  `register_no` varchar(50) NOT NULL,
  `name` varchar(100) NOT NULL,
  `password` varchar(255) NOT NULL,
  `department` varchar(50) NOT NULL,
  `email` varchar(100) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `hods`
--

INSERT INTO `hods` (`id`, `register_no`, `name`, `password`, `department`, `email`) VALUES
(1, '1204', 'Vijayalakshmi K', 'c55b3f6a81a515bffa332d727560f1b5', 'CSE', 'anishrkumar2k5@gmail.com');

-- --------------------------------------------------------

--
-- Table structure for table `jas`
--

CREATE TABLE `jas` (
  `id` int(11) NOT NULL,
  `register_no` varchar(50) NOT NULL,
  `name` varchar(100) NOT NULL,
  `password` varchar(255) NOT NULL,
  `department` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `jas`
--

INSERT INTO `jas` (`id`, `register_no`, `name`, `password`, `department`) VALUES
(2, '2', 'Rishi', 'c55b3f6a81a515bffa332d727560f1b5', 'CSE'),
(5, '2322', 'Testing JA Finalist', 'c55b3f6a81a515bffa332d727560f1b5', 'CSE'),
(6, '2407', 'JA Testing', 'c55b3f6a81a515bffa332d727560f1b5', 'CSE');

-- --------------------------------------------------------

--
-- Table structure for table `lab_form`
--

CREATE TABLE `lab_form` (
  `id` int(11) NOT NULL,
  `registerNumber` varchar(20) NOT NULL,
  `studentName` varchar(100) NOT NULL,
  `year` varchar(20) NOT NULL,
  `department` varchar(50) NOT NULL,
  `section` varchar(10) NOT NULL,
  `purpose` enum('internal','external') NOT NULL,
  `fullDayDate` date DEFAULT NULL,
  `fromTime` time DEFAULT NULL,
  `toTime` time DEFAULT NULL,
  `fromDate` date DEFAULT NULL,
  `toDate` date DEFAULT NULL,
  `collegeName` varchar(150) DEFAULT NULL,
  `eventName` varchar(150) DEFAULT NULL,
  `extFromDate` date DEFAULT NULL,
  `extToDate` date DEFAULT NULL,
  `mentor` varchar(100) NOT NULL,
  `systemRequired` enum('Yes','No') NOT NULL,
  `submitted_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `lab_technicians`
--

CREATE TABLE `lab_technicians` (
  `id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `department` varchar(50) DEFAULT NULL,
  `register_no` varchar(20) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `email` varchar(100) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `lab_technicians`
--

INSERT INTO `lab_technicians` (`id`, `name`, `department`, `register_no`, `password`, `created_at`, `email`) VALUES
(1, 'Suresh  Kumar', 'CSE', '2030', 'c55b3f6a81a515bffa332d727560f1b5', '2025-10-11 05:35:46', '953623104007@ritrjpm.ac.in');

-- --------------------------------------------------------

--
-- Table structure for table `mentors`
--

CREATE TABLE `mentors` (
  `id` int(11) NOT NULL,
  `register_no` varchar(50) NOT NULL,
  `name` varchar(100) NOT NULL,
  `password` varchar(255) NOT NULL,
  `department` varchar(50) NOT NULL,
  `year` int(11) DEFAULT NULL,
  `section` varchar(10) NOT NULL,
  `mentor_email` varchar(100) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `mentors`
--

INSERT INTO `mentors` (`id`, `register_no`, `name`, `password`, `department`, `year`, `section`, `mentor_email`) VALUES
(10, '1232', 'Vivek V', 'c55b3f6a81a515bffa332d727560f1b5', 'CSE', 3, 'A', 'davidw.joshua@gmail.com'),
(15, '1230', 'Swarna Sudha', 'c55b3f6a81a515bffa332d727560f1b5', 'CSE', 3, 'A', '953623104007@ritrjpm.ac.in'),
(17, '1289', 'Vijay Amala Devi', 'c55b3f6a81a515bffa332d727560f1b5', 'CSE', 3, 'B', 'anishkumarfrom2k5@gmail.com');

-- --------------------------------------------------------

--
-- Table structure for table `od_applications`
--

CREATE TABLE `od_applications` (
  `id` int(11) NOT NULL,
  `register_no` varchar(20) DEFAULT NULL,
  `student_name` varchar(100) DEFAULT NULL,
  `year` varchar(20) DEFAULT NULL,
  `department` varchar(50) DEFAULT NULL,
  `section` varchar(10) DEFAULT NULL,
  `mentor` varchar(100) DEFAULT NULL,
  `purpose` text DEFAULT NULL,
  `od_type` varchar(20) DEFAULT NULL,
  `od_date` date DEFAULT NULL,
  `from_time` time DEFAULT NULL,
  `to_time` time DEFAULT NULL,
  `from_date` date DEFAULT NULL,
  `to_date` date DEFAULT NULL,
  `college_name` varchar(150) DEFAULT NULL,
  `event_name` varchar(150) DEFAULT NULL,
  `lab_required` tinyint(1) DEFAULT 0,
  `lab_name` varchar(255) DEFAULT NULL,
  `system_required` tinyint(1) DEFAULT 0,
  `request_bonafide` tinyint(1) DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `status` varchar(20) DEFAULT 'Pending'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `od_applications`
--

INSERT INTO `od_applications` (`id`, `register_no`, `student_name`, `year`, `department`, `section`, `mentor`, `purpose`, `od_type`, `od_date`, `from_time`, `to_time`, `from_date`, `to_date`, `college_name`, `event_name`, `lab_required`, `lab_name`, `system_required`, `request_bonafide`, `created_at`, `status`) VALUES
(57, '953623104053', 'Lalith Krishna V M', '3rd Year', 'CSE', 'A', 'Vivek V', 'External Hackathon @ VIT', 'external', '0000-00-00', '00:00:00', '00:00:00', '2025-10-08', '2025-10-10', 'Vellore Institute Of Technology', 'Vellorathon', 0, NULL, 0, 1, '2025-10-08 03:54:05', 'Mentors Rejected'),
(58, '953623104053', 'Lalith Krishna V M', '3rd Year', 'CSE', 'A', 'Vivek V', 'Internal Hackathon @ RIT', 'internal', '2025-10-08', '10:27:00', '11:27:00', '2025-10-08', '2025-10-08', '', '', 0, NULL, 0, 0, '2025-10-08 04:57:56', 'HOD Accepted'),
(59, '953623104044', 'S Jeyaseelan', '3rd Year', 'CSE', 'A', 'Vivek V', 'KPR Institute Of Technology', 'internal', '2025-10-11', '00:00:00', '00:00:00', '2025-10-11', '2025-10-11', 'Kpr ', 'Hackaxerlate', 0, NULL, 0, 0, '2025-10-11 04:24:11', 'HOD Rejected'),
(70, '953623104053', 'Lalith Krishna V M', '3rd Year', 'CSE', 'A', 'Vivek V', '2K6 Hackathon', 'internal', '0000-00-00', '10:53:00', '12:53:00', '2025-10-11', '2025-10-11', '', '', 1, 'IoT Lab', 1, 0, '2025-10-11 05:21:46', 'HOD Accepted'),
(71, '953623104053', 'Lalith Krishna V M', '3rd Year', 'CSE', 'A', 'Vivek V', 'PEC Hacks!!', 'internal', '2005-11-20', '00:00:00', '00:00:00', '2005-11-20', '2005-11-20', '', '', 1, 'Network Lab', 0, 0, '2025-10-13 03:52:56', 'HOD Accepted'),
(72, '953623104007', 'Anish Kumar', '3rd Year', 'CSE', 'A', 'Vivek V', 'Going for Hackathon', 'internal', '0000-00-00', '12:00:00', '17:30:00', '2025-10-14', '2025-10-14', '', '', 0, NULL, 0, 0, '2025-10-14 05:21:28', 'Pending'),
(73, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing Period Wise thing', 'internal', '0000-00-00', '12:00:00', '15:00:00', '2025-10-14', '2025-10-14', '', '', 0, NULL, 0, 0, '2025-10-14 05:24:51', 'HOD Rejected'),
(74, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Checking For OD_date for Period Wise', 'internal', '0000-00-00', '15:00:00', '16:00:00', '2025-10-14', '2025-10-14', '', '', 0, NULL, 0, 0, '2025-10-14 05:29:27', 'HOD Accepted'),
(75, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Final Testing Whether it is asking Lab or Not and i am giving no lab', 'internal', '0000-00-00', '09:00:00', '14:00:00', '2025-10-14', '2025-10-14', '', '', 0, NULL, 0, 0, '2025-10-14 06:48:04', 'HOD Accepted'),
(76, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'No Lab', 'internal', '0000-00-00', '13:00:00', '15:00:00', '2025-10-14', '2025-10-14', '', '', 0, NULL, 0, 0, '2025-10-14 06:55:45', 'HOD Accepted'),
(77, '953623104007', 'Anish Kumar', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing the PHPMailer', 'internal', '2025-10-16', '00:00:00', '00:00:00', '2025-10-16', '2025-10-16', '', '', 0, NULL, 0, 0, '2025-10-15 04:39:05', 'Pending'),
(78, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Going to the  PHPMailer ', 'internal', '0000-00-00', '12:00:00', '14:00:00', '2025-10-15', '2025-10-15', '', '', 0, NULL, 0, 0, '2025-10-15 08:44:24', 'Mentors Rejected'),
(79, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Going for Testing ', 'internal', '0000-00-00', '12:00:00', '14:00:00', '2025-10-15', '2025-10-15', '', '', 0, NULL, 0, 0, '2025-10-15 09:04:37', 'HOD Accepted'),
(80, '953623104066', 'K Mythri Vaishnavi', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Raja Nursery', 'internal', '0000-00-00', '00:00:00', '00:00:00', '2025-08-09', '2052-05-09', '', '', 0, NULL, 0, 0, '2025-10-15 09:40:15', 'HOD Accepted'),
(81, '953623104007', 'Anish Kumar', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'This is for Testing the PHP Mailer it receicve my mail or not', 'internal', '2005-11-20', '00:00:00', '00:00:00', '2005-11-20', '2005-11-20', '', '', 0, NULL, 0, 0, '2025-10-15 12:13:28', 'Pending'),
(82, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Second Testing', 'internal', '0200-11-20', '00:00:00', '00:00:00', '0200-11-20', '0200-11-20', '', '', 0, NULL, 0, 0, '2025-10-15 12:22:26', 'Pending'),
(83, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Purpose of using phpmailer', 'internal', '2005-11-20', '00:00:00', '00:00:00', '2005-11-20', '2005-11-20', '', '', 0, NULL, 0, 0, '2025-10-16 01:21:24', 'Pending'),
(84, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Normal file mail testing', 'internal', '2005-10-30', '00:00:00', '00:00:00', '2005-10-30', '2005-10-30', '', '', 0, NULL, 0, 0, '2025-10-16 01:25:10', 'Pending'),
(85, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'testing', 'internal', '2005-11-25', '00:00:00', '00:00:00', '2005-11-25', '2005-11-25', '', '', 0, NULL, 0, 0, '2025-10-16 01:41:20', 'Pending'),
(86, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Tester', 'internal', '2005-11-20', '00:00:00', '00:00:00', '2005-11-20', '2005-11-20', '', '', 0, NULL, 0, 0, '2025-10-16 01:44:06', 'Pending'),
(87, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing Mailer', 'internal', '2005-10-29', '00:00:00', '00:00:00', '2005-10-29', '2005-10-29', '', '', 0, NULL, 0, 0, '2025-10-16 02:13:51', 'Pending'),
(88, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing with Mailer', 'internal', '2053-06-27', '00:00:00', '00:00:00', '2053-06-27', '2053-06-27', '', '', 0, NULL, 0, 0, '2025-10-16 03:49:17', 'Pending'),
(89, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Test PHPMailer thing', 'internal', '2003-11-20', '00:00:00', '00:00:00', '2003-11-20', '2003-11-20', '', '', 0, NULL, 0, 0, '2025-10-16 03:55:14', 'Pending'),
(90, '953623104066', 'test', '3rd Year', 'CSE', 'B', 'Devi', 'dwdwwwew', 'internal', '2025-10-16', '00:00:00', '00:00:00', '2025-10-16', '2025-10-16', '', '', 0, NULL, 0, 0, '2025-10-16 03:57:14', 'Pending'),
(91, '953623104066', 'Anish Kumar R', '3rd Year', 'CSE', 'B', 'Devi', 'xduiwndfpefjefe', 'internal', '2025-10-16', '00:00:00', '00:00:00', '2025-10-16', '2025-10-16', '', '', 0, NULL, 0, 0, '2025-10-16 04:24:14', 'Pending'),
(92, '953623104066', 'Anish Kumar R', '3rd Year', 'CSE', 'B', 'Devi', 'ertfyguhijpemfw fhvweghfje', 'internal', '2025-10-16', '00:00:00', '00:00:00', '2025-10-16', '2025-10-16', '', '', 0, NULL, 0, 0, '2025-10-16 04:26:30', 'Pending'),
(93, '953623104066', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'zrdxtcfyvgubhinjok', 'internal', '2025-10-16', '00:00:00', '00:00:00', '2025-10-16', '2025-10-16', '', '', 0, NULL, 0, 0, '2025-10-16 04:41:15', 'Pending'),
(94, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Just testing ', 'internal', '2006-11-23', '00:00:00', '00:00:00', '2006-11-23', '2006-11-23', '', '', 0, NULL, 0, 0, '2025-10-16 05:03:37', 'Pending'),
(95, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Jest Testing', 'internal', '2510-05-04', '00:00:00', '00:00:00', '2510-05-04', '2510-05-04', '', '', 0, NULL, 0, 0, '2025-10-16 05:42:02', 'Pending'),
(96, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'This is PHP Tester', 'internal', '2342-11-23', '00:00:00', '00:00:00', '2342-11-23', '2342-11-23', '', '', 0, NULL, 0, 0, '2025-10-16 05:51:14', 'Pending'),
(97, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Nothing Just thing about PHP Mailer', 'internal', '2025-12-17', '00:00:00', '00:00:00', '2025-12-17', '2025-12-17', '', '', 0, NULL, 0, 0, '2025-10-16 05:57:31', 'Pending'),
(98, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Test With Real Time Data', 'internal', '2025-11-20', '00:00:00', '00:00:00', '2025-11-20', '2025-11-20', '', '', 0, NULL, 0, 0, '2025-10-16 06:37:00', 'Pending'),
(99, '953623104007', 'Mythri Vaishnavi', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'poriyaan', 'internal', '6325-05-07', '00:00:00', '00:00:00', '6325-05-07', '6325-05-07', '', '', 0, NULL, 0, 0, '2025-10-16 09:36:01', 'Mentors Rejected'),
(100, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Wested thing', 'internal', '2002-11-19', '00:00:00', '00:00:00', '2002-11-19', '2002-11-19', '', '', 0, NULL, 0, 0, '2025-10-16 10:31:04', 'Pending'),
(101, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'This is Purpose thing i had changed many thing', 'internal', '2005-12-21', '00:00:00', '00:00:00', '2005-12-21', '2005-12-21', '', '', 0, NULL, 0, 0, '2025-10-17 14:34:59', 'Mentors Accepted'),
(102, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Test it', 'internal', '0000-00-00', '14:00:00', '16:00:00', '2025-10-17', '2025-10-17', '', '', 0, NULL, 0, 0, '2025-10-17 14:39:44', 'Mentors Accepted'),
(103, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing it things', 'internal', '0000-00-00', '00:00:00', '00:00:00', '2025-10-17', '2025-10-21', '', '', 0, NULL, 0, 0, '2025-10-17 15:16:47', 'Mentors Accepted'),
(104, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing 2 nothing but that i hope this will come', 'internal', '2005-11-24', '00:00:00', '00:00:00', '2005-11-24', '2005-11-24', '', '', 0, NULL, 0, 0, '2025-10-17 15:28:15', 'Pending'),
(105, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Last Testing here', 'internal', '2025-11-30', '00:00:00', '00:00:00', '2025-11-30', '2025-11-30', '', '', 0, NULL, 0, 0, '2025-10-17 15:32:16', 'Mentors Accepted'),
(106, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Purpose of testing', 'internal', '2005-11-20', '00:00:00', '00:00:00', '2005-11-20', '2005-11-20', '', '', 0, NULL, 0, 0, '2025-10-22 02:05:07', 'Mentors Rejected'),
(107, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'nOHTING JUST FUNCTIONING THE TEST MAILER', 'internal', '2005-11-20', '00:00:00', '00:00:00', '2005-11-20', '2005-11-20', '', '', 1, 'Hardware Lab', 0, 0, '2025-10-22 02:08:56', 'Mentors Rejected'),
(108, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Zero Testing', 'internal', '2025-11-20', '00:00:00', '00:00:00', '2025-11-20', '2025-11-20', '', '', 0, NULL, 0, 0, '2025-10-22 03:52:17', 'Mentors Rejected'),
(109, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing at 22-10-2025 (9:34)', 'internal', '2025-10-22', '00:00:00', '00:00:00', '2025-10-22', '2025-10-22', '', '', 1, 'IoT Lab', 0, 0, '2025-10-22 04:04:31', 'Mentors Rejected'),
(110, '953623104042', 'Saravanankamalesh', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'skjfsdnfksf', 'internal', '2252-11-20', '00:00:00', '00:00:00', '2252-11-20', '2252-11-20', '', '', 1, 'IoT Lab', 0, 0, '2025-10-22 04:06:58', 'Mentors Rejected'),
(111, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing at 22-10-2025', 'internal', '2925-10-29', '00:00:00', '00:00:00', '2925-10-29', '2925-10-29', '', '', 1, 'IoT Lab', 0, 0, '2025-10-22 04:18:12', 'Mentors Rejected'),
(112, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'TESTING 4 STAGE', 'internal', '2025-11-25', '00:00:00', '00:00:00', '2025-11-25', '2025-11-25', '', '', 1, 'Network Lab', 0, 0, '2025-10-22 04:19:49', 'Mentors Rejected'),
(113, '953623104066', 'Mythri Vaishnavi', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing with MV', 'internal', '2025-10-23', '00:00:00', '00:00:00', '2025-10-23', '2025-10-23', '', '', 1, 'Data Science Lab', 1, 0, '2025-10-22 04:33:57', 'Mentors Rejected'),
(114, '953623104066', 'Mythri Vaishnavi', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing with MV', 'internal', '2025-10-23', '00:00:00', '00:00:00', '2025-10-23', '2025-10-23', '', '', 1, 'Data Science Lab', 1, 0, '2025-10-22 04:34:11', 'Mentors Rejected'),
(115, '953623104066', 'Mythri Vaishnavi', '3rd Year', 'CSE', 'B', 'Devi', 'I am Going for Hackathon in Ramco Mills', 'internal', '2025-10-24', '00:00:00', '00:00:00', '2025-10-24', '2025-10-24', '', '', 1, 'Network Lab', 1, 0, '2025-10-22 04:38:52', 'Pending'),
(116, '953623104066', 'Mythri Vaishnavi', '3rd Year', 'CSE', 'A', 'Vivek V', 'Final Testing for the Mail from Student -> Mentor', 'internal', '2025-10-23', '00:00:00', '00:00:00', '2025-10-23', '2025-10-23', '', '', 1, 'IoT Lab', 0, 0, '2025-10-22 04:44:14', 'Pending'),
(117, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Final Testing from Student -> Mentor', 'internal', '2025-10-23', '00:00:00', '00:00:00', '2025-10-23', '2025-10-23', '', '', 1, 'AI Lab', 0, 0, '2025-10-22 04:46:45', 'Mentors Rejected'),
(118, '953623104066', 'Mythri Vaishnavi', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Final Testing with team Member', 'internal', '2025-10-24', '00:00:00', '00:00:00', '2025-10-24', '2025-10-24', '', '', 1, 'IoT Lab', 0, 0, '2025-10-22 04:50:15', 'Mentors Accepted'),
(119, '953623104007', 'Anish KuA', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing for Team Members', 'internal', '2925-11-20', '00:00:00', '00:00:00', '2925-11-20', '2925-11-20', '', '', 1, 'Network Lab', 0, 0, '2025-10-22 05:14:19', 'Mentors Rejected'),
(120, '953623104007', 'Anish Kumar ', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'This is Testing Final Testing for Team Member i think it is correct', 'internal', '2025-12-23', '00:00:00', '00:00:00', '2025-12-23', '2025-12-23', '', '', 1, 'Hardware Lab', 0, 0, '2025-10-22 05:31:26', 'HOD Accepted'),
(121, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Nothing more about it just testing', 'internal', '2005-10-22', '00:00:00', '00:00:00', '2005-10-22', '2005-10-22', '', '', 1, 'Hardware Lab', 0, 0, '2025-10-22 09:47:53', 'HOD Accepted'),
(122, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Vivek V', 'This is for testing and give me the change of this', 'internal', '2025-10-23', '00:00:00', '00:00:00', '2025-10-23', '2025-10-23', '', '', 0, NULL, 0, 0, '2025-10-22 12:09:30', 'Mentors Rejected'),
(123, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Test Case', 'internal', '0000-00-00', '00:00:00', '00:00:00', '0000-00-00', '0000-00-00', '', '', 0, NULL, 0, 0, '2025-10-23 04:30:22', 'Mentors Rejected'),
(124, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing thing about HOD from mentor mail', 'internal', '2025-10-23', '00:00:00', '00:00:00', '2025-10-23', '2025-10-23', '', '', 1, 'Network Lab', 1, 0, '2025-10-23 06:11:30', 'Mentors Accepted'),
(125, '953623104013', 'Diliban S', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Latest testing', 'internal', '2025-10-23', '00:00:00', '00:00:00', '2025-10-23', '2025-10-23', '', '', 1, 'AI Lab', 0, 0, '2025-10-23 06:27:03', 'Mentors Accepted'),
(126, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Mentor to HOD Email Notificaiton', 'internal', '2005-10-25', '00:00:00', '00:00:00', '2005-10-25', '2005-10-25', '', '', 0, NULL, 0, 0, '2025-10-23 06:46:14', 'Mentors Accepted'),
(127, '953623104007', 'Aish Kumar', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'nothing', 'internal', '2025-11-20', '00:00:00', '00:00:00', '2025-11-20', '2025-11-20', '', '', 0, NULL, 0, 0, '2025-10-23 06:51:48', 'Mentors Accepted'),
(128, '953623104007', 'Anish kUmar ', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Nothing much about it', 'internal', '2025-10-25', '00:00:00', '00:00:00', '2025-10-25', '2025-10-25', '', '', 0, NULL, 0, 0, '2025-10-23 06:54:17', 'Mentors Accepted'),
(129, '953623104007', 'Anish Kumar  R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Latest testing at 25.10.2025', 'internal', '2025-10-26', '00:00:00', '00:00:00', '2025-10-26', '2025-10-26', '', '', 0, NULL, 0, 0, '2025-10-25 04:12:55', 'Mentors Accepted'),
(130, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Lastst Testing at 9:47', 'internal', '2025-10-27', '00:00:00', '00:00:00', '2025-10-27', '2025-10-27', '', '', 0, NULL, 0, 0, '2025-10-25 04:17:59', 'Pending'),
(131, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing at 9:49', 'internal', '2025-10-26', '00:00:00', '00:00:00', '2025-10-26', '2025-10-26', '', '', 0, NULL, 0, 0, '2025-10-25 04:20:11', 'Mentors Accepted'),
(132, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing at 9:56', 'internal', '2025-10-26', '00:00:00', '00:00:00', '2025-10-26', '2025-10-26', '', '', 0, NULL, 0, 0, '2025-10-25 04:25:42', 'Mentors Accepted'),
(133, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'This is the purpose of the testing the email sending from the mentor to the hod', 'internal', '2005-10-26', '00:00:00', '00:00:00', '2005-10-26', '2005-10-26', '', '', 0, NULL, 0, 0, '2025-10-25 04:30:42', 'HOD Accepted'),
(134, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing @ 10:18', 'internal', '2025-10-27', '00:00:00', '00:00:00', '2025-10-27', '2025-10-27', '', '', 1, 'AI Lab', 0, 0, '2025-10-25 04:49:00', 'HOD Rejected'),
(135, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing @ 10:22', 'internal', '2025-10-28', '00:00:00', '00:00:00', '2025-10-28', '2025-10-28', '', '', 1, 'Data Science Lab', 1, 0, '2025-10-25 04:52:55', 'HOD Rejected'),
(136, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing @ 10:30', 'internal', '2025-10-26', '00:00:00', '00:00:00', '2025-10-26', '2025-10-26', '', '', 1, 'Network Lab', 0, 0, '2025-10-25 05:00:39', 'HOD Accepted'),
(137, '953623104049', 'Kathir Vel', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing @ 10:40', 'internal', '2025-10-25', '00:00:00', '00:00:00', '2025-10-25', '2025-10-25', '', '', 1, 'Hardware Lab', 0, 0, '2025-10-25 05:11:13', 'Mentors Accepted'),
(138, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'This is for external prinicpal email testing ', 'external', '0000-00-00', '00:00:00', '00:00:00', '2025-10-26', '2025-10-27', 'Easwari Engineering College', '30 Hours Hackathon', 0, NULL, 0, 1, '2025-10-25 05:29:31', 'HOD Accepted'),
(139, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'This is 1st testing after kolarubodi', 'internal', '2005-11-20', '00:00:00', '00:00:00', '2005-11-20', '2005-11-20', '', '', 0, NULL, 0, 0, '2025-10-27 07:33:23', 'Mentors Accepted'),
(140, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Purpose of testing thing ', 'internal', '2025-10-24', '00:00:00', '00:00:00', '2025-10-24', '2025-10-24', '', '', 1, 'IoT Lab', 0, 0, '2025-10-30 02:15:39', 'Pending'),
(141, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Purpose of testing thing ', 'internal', '2025-10-24', '00:00:00', '00:00:00', '2025-10-24', '2025-10-24', '', '', 1, 'IoT Lab', 0, 0, '2025-10-30 02:15:45', 'Mentors Accepted'),
(142, '953623104007', 'Anish Kumar ', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Purpose of testing the email', 'internal', '2005-10-31', '00:00:00', '00:00:00', '2005-10-31', '2005-10-31', '', '', 1, 'Network Lab', 0, 0, '2025-10-30 02:17:31', 'Mentors Accepted'),
(143, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Final Testing ', 'internal', '2025-10-28', '00:00:00', '00:00:00', '2025-10-28', '2025-10-28', '', '', 0, NULL, 0, 0, '2025-10-30 02:32:12', 'HOD Accepted'),
(144, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'This is for team ', 'internal', '2025-11-29', '00:00:00', '00:00:00', '2025-11-29', '2025-11-29', '', '', 1, 'Network Lab', 0, 0, '2025-10-31 08:51:15', 'HOD Accepted'),
(145, '953623104007', 'Anish kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Final Testing ', 'internal', '2025-11-28', '00:00:00', '00:00:00', '2025-11-28', '2025-11-28', '', '', 1, 'IoT Lab', 0, 0, '2025-10-31 09:02:20', 'HOD Accepted'),
(146, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'This  is for testing for the following testing tools ', 'internal', '2025-11-29', '00:00:00', '00:00:00', '2025-11-29', '2025-11-29', '', '', 1, 'Network Lab', 0, 0, '2025-10-31 13:53:52', 'HOD Accepted'),
(147, '953623104007', 'AK', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'This is enabled lab requirement', 'internal', '2025-11-02', '00:00:00', '00:00:00', '2025-11-02', '2025-11-02', '', '', 1, NULL, 0, 0, '2025-10-31 14:05:35', 'HOD Accepted'),
(148, '953623104007', 'aNISH kUMAR ', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Final Testing', 'internal', '2025-11-11', '00:00:00', '00:00:00', '2025-11-11', '2025-11-11', '', '', 1, 'AI Lab', 0, 0, '2025-10-31 14:09:21', 'HOD Accepted'),
(149, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'This is first Final Testing', 'internal', '2025-11-02', '00:00:00', '00:00:00', '2025-11-02', '2025-11-02', '', '', 0, NULL, 0, 0, '2025-11-01 02:01:45', 'HOD Accepted'),
(150, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'This is for External', 'external', '0000-00-00', '00:00:00', '00:00:00', '2025-11-21', '2025-11-22', 'Easwari Engineering College', '30 Hours Hackathon', 0, NULL, 0, 0, '2025-11-01 02:07:00', 'Principal Accepted'),
(151, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'This is for Team Internal Testing', 'internal', '0000-00-00', '09:43:00', '16:43:00', '2025-11-01', '2025-11-01', '', '', 1, 'AI Lab', 0, 0, '2025-11-01 02:13:28', 'HOD Accepted'),
(152, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'This is for External ', 'external', '0000-00-00', '00:00:00', '00:00:00', '2025-11-01', '2025-11-02', 'St.Joseph Engineering College', '30 Hours Hackathon', 0, NULL, 0, 0, '2025-11-01 02:24:13', 'Principal Accepted'),
(153, '953623104007', 'Anish Kumatr R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Testing final for is this ', 'internal', '0000-00-00', '02:45:00', '03:40:00', '2025-11-04', '2025-11-04', '', '', 0, NULL, 0, 0, '2025-11-04 10:07:39', 'Pending'),
(154, '953623104007', 'Anish Kumar', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'This sis for harish  purpose', 'internal', '2025-11-06', '00:00:00', '00:00:00', '2025-11-06', '2025-11-06', '', '', 1, 'Network Lab', 0, 0, '2025-11-04 10:27:59', 'HOD Accepted'),
(155, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Final Testing of OD-Module', 'internal', '2025-12-16', '00:00:00', '00:00:00', '2025-12-16', '2025-12-16', '', '', 1, 'IoT Lab', 0, 0, '2025-12-16 02:27:53', 'HOD Accepted'),
(156, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'This is for Official Testing', 'internal', '0000-00-00', '00:00:00', '00:00:00', '2025-12-17', '2025-12-18', '', '', 1, 'AI Lab', 0, 0, '2025-12-16 07:51:25', 'Mentors Accepted'),
(157, '953623104007', 'Anish Kumar R', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'tESTING', 'internal', '0000-00-00', '00:00:00', '00:00:00', '2025-12-16', '2025-12-18', '', '', 1, 'Network Lab', 0, 0, '2025-12-16 07:54:53', 'HOD Accepted'),
(158, '953623104007', 'Anish Kumar R', '3', 'CSE', 'A', 'Swarna Sudha', 'OD-Module Final Stage Testing', 'internal', '2026-01-22', '00:00:00', '00:00:00', '2026-01-22', '2026-01-22', '', '', 1, 'AI Lab', 0, 0, '2026-01-21 09:50:14', 'Pending'),
(159, '953623104007', 'Anish Kumar R', '3', 'CSE', 'A', 'Swarna Sudha', 'ON DUTY EMAIL TESTING', 'internal', '2026-01-23', '00:00:00', '00:00:00', '2026-01-23', '2026-01-23', '', '', 1, 'AI Lab', 0, 0, '2026-01-21 09:53:30', 'Mentors Accepted'),
(160, '953623104007', 'Anish Kumar R', '3', 'CSE', 'A', 'Swarna Sudha', 'OD', 'internal', '2026-12-26', '00:00:00', '00:00:00', '2026-12-26', '2026-12-26', '', '', 0, NULL, 0, 0, '2026-01-21 10:22:58', 'Pending'),
(161, '953623104007', 'Anish Kumar R', '3', 'CSE', 'A', 'Swarna Sudha', 'OD Testing Final Stage', 'internal', '2026-01-23', '00:00:00', '00:00:00', '2026-01-23', '2026-01-23', '', '', 1, 'AI Lab', 0, 0, '2026-01-22 07:14:11', 'Pending'),
(162, '953623104007', 'Anish Kumar R', '3', 'CSE', 'A', 'Swarna Sudha', 'OD Final Testing', 'internal', '2026-11-23', '00:00:00', '00:00:00', '2026-11-23', '2026-11-23', '', '', 1, 'AI Lab', 0, 0, '2026-01-22 07:16:12', 'HOD Accepted'),
(163, '953623104007', 'Anish Kumar R', '3', 'CSE', 'A', 'Swarna Sudha', 'OnDuty', 'internal', '2026-01-23', '00:00:00', '00:00:00', '2026-01-23', '2026-01-23', '', '', 1, 'AI Lab', 0, 0, '2026-01-22 07:26:18', 'Pending'),
(164, '953623104007', 'Anish Kumar R', '3', 'CSE', 'A', 'Swarna Sudha', 'Testing On Duty Final Stage Test', 'internal', '2026-01-23', '00:00:00', '00:00:00', '2026-01-23', '2026-01-23', '', '', 1, 'AI Lab', 0, 0, '2026-01-22 07:39:02', 'Pending'),
(165, '953623104007', 'Anish Kumar R', '3', 'CSE', 'A', 'Swarna Sudha', 'On Duty Final Testing Preparation Last Stage and Final Stage', 'internal', '2026-01-24', '00:00:00', '00:00:00', '2026-01-24', '2026-01-24', '', '', 1, 'Data Science Lab', 0, 0, '2026-01-22 08:36:22', 'HOD Accepted'),
(166, '953623104021', 'David Joshua', '3', 'CSE', 'B', 'Vijay Amala Devi', 'On Duty  For  External', 'external', '0000-00-00', '00:00:00', '00:00:00', '2026-01-22', '2026-01-26', 'Easwari Engineering College', '30 Hours Hackathon', 0, NULL, 0, 1, '2026-01-22 08:43:10', 'HOD Accepted'),
(167, '953623104007', 'Anish Kumar R', '3', 'CSE', 'A', 'Swarna Sudha', 'On Duty Model Training', 'external', '0000-00-00', '00:00:00', '00:00:00', '2026-01-23', '2026-01-25', 'Easwari Engineering College', '30 Hours Hackathon', 0, NULL, 0, 1, '2026-01-22 09:12:38', 'Pending'),
(168, '953623104007', 'Anish Kumar R', '3', 'CSE', 'A', 'Swarna Sudha', 'Duty', 'internal', '2026-01-24', '00:00:00', '00:00:00', '2026-01-24', '2026-01-24', '', '', 1, 'IoT Lab', 0, 0, '2026-01-22 09:18:28', 'HOD Rejected'),
(169, '953623104007', 'Anish Kumar R', '3', 'CSE', 'A', 'Swarna Sudha', 'Only On Duty Thing', 'internal', '2026-11-24', '00:00:00', '00:00:00', '2026-11-24', '2026-11-24', '', '', 1, 'AI Lab', 0, 0, '2026-01-22 09:24:17', 'Mentors Rejected'),
(170, '953623104053', 'Lalith Krishna V M', '3', 'CSE', 'A', 'Swarna Sudha', 'Testing', 'internal', '2026-01-23', '00:00:00', '00:00:00', '2026-01-23', '2026-01-23', '', '', 0, NULL, 0, 0, '2026-01-22 09:26:56', 'HOD Accepted'),
(171, '953623104007', 'Anish Kumar R', '3', 'CSE', 'A', 'Swarna Sudha', 'on duty final tresting ', 'internal', '2026-01-24', '00:00:00', '00:00:00', '2026-01-24', '2026-01-24', '', '', 0, NULL, 0, 0, '2026-01-22 10:59:59', 'HOD Rejected'),
(172, '953623104021', 'David W Joshua', '3', 'CSE', 'B', 'Vijay Amala Devi', 'This is for testing', 'internal', '2026-01-24', '00:00:00', '00:00:00', '2026-01-24', '2026-01-24', '', '', 1, 'AI Lab', 0, 0, '2026-01-23 05:40:10', 'HOD Accepted'),
(173, '953623104003', 'Afridha Fathima Syed Ibrahim', '3', 'CSE', 'B', 'Vijay Amala Devi', 'Project', 'internal', '0000-00-00', '00:00:00', '00:00:00', '2026-01-23', '2026-01-28', '', '', 1, 'AI Lab', 1, 0, '2026-01-23 05:49:29', 'HOD Accepted'),
(174, '953623104021', 'David W Joshua', '3', 'CSE', 'B', 'Vijay Amala Devi', 'On Duty', 'internal', '2026-01-24', '00:00:00', '00:00:00', '2026-01-24', '2026-01-24', '', '', 1, 'AI Lab', 0, 0, '2026-01-23 10:31:03', 'Pending'),
(175, '953623104040', 'Janaramji R', '3', 'CSE', 'A', 'Vivek V', 'Selected for Hackathon @ sri eshwar college of engineering', 'internal', '0000-00-00', '12:09:00', '14:10:00', '2026-01-24', '2026-01-24', '', '', 1, 'IoT Lab', 0, 0, '2026-01-24 05:42:01', 'Pending'),
(176, '953623104021', 'David W Joshua', '3', 'CSE', 'B', 'Vijay Amala Devi', 'mgbjabg', 'internal', '2026-02-20', '00:00:00', '00:00:00', '2026-02-20', '2026-02-20', '', '', 1, 'Open Source Lab', 0, 0, '2026-02-03 09:48:21', 'Pending'),
(177, '953623104021', 'David W Joshua', '3', 'CSE', 'B', 'Vijay Amala Devi', ',jfg', 'internal', '0000-00-00', '14:00:00', '15:00:00', '2026-02-13', '2026-02-13', '', '', 1, 'IOS Lab', 0, 0, '2026-02-13 09:38:24', 'Pending'),
(178, '953623104007', 'Anish Kumar R', '3', 'CSE', 'A', NULL, 'This is for Testing Final before Deployment', 'internal', '2026-03-11', NULL, NULL, NULL, NULL, '', '', 1, 'Open Source Lab', 0, 0, '2026-03-11 04:28:21', 'HOD Accepted');

-- --------------------------------------------------------

--
-- Table structure for table `od_team_members`
--

CREATE TABLE `od_team_members` (
  `id` int(11) NOT NULL,
  `od_id` int(11) DEFAULT NULL,
  `member_name` varchar(100) DEFAULT NULL,
  `member_regno` varchar(20) DEFAULT NULL,
  `member_year` varchar(20) DEFAULT NULL,
  `member_department` varchar(50) DEFAULT NULL,
  `member_section` varchar(10) DEFAULT NULL,
  `mentor` varchar(255) DEFAULT NULL,
  `mentor_status` enum('Pending','Accepted','Rejected') DEFAULT 'Pending'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `od_team_members`
--

INSERT INTO `od_team_members` (`id`, `od_id`, `member_name`, `member_regno`, `member_year`, `member_department`, `member_section`, `mentor`, `mentor_status`) VALUES
(65, 57, 'Lalith Krishna V M', '953623104053', '3rd Year', 'CSE', 'A', 'Vivek V', 'Accepted'),
(66, 57, 'Anish Kumar', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Rejected'),
(67, 57, 'Diliban S M', '953623104013', '3rd Year', 'CSE', 'A', 'Vivek V', 'Accepted'),
(68, 58, 'Lalith Krishna V M', '953623104053', '3rd Year', 'CSE', 'A', 'Vivek V', 'Accepted'),
(69, 58, 'Anish Kumar', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Rejected'),
(70, 58, 'Diliban S M', '953623104013', '3rd Year', 'CSE', 'A', 'Vivek V', 'Accepted'),
(71, 59, 'S Jeyaseelan', '953623104044', '3rd Year', 'CSE', 'A', 'Vivek V', 'Accepted'),
(90, 70, 'Lalith Krishna V M', '953623104053', '3rd Year', 'CSE', 'A', 'Vivek V', 'Accepted'),
(91, 70, 'Jeyseelan S', '953623104044', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(92, 71, 'Lalith Krishna V M', '953623104053', '3rd Year', 'CSE', 'A', 'Vivek V', 'Accepted'),
(93, 71, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(94, 72, 'Anish Kumar', '953623104007', '3rd Year', 'CSE', 'A', 'Vivek V', 'Pending'),
(95, 73, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(96, 74, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(97, 75, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(98, 76, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(99, 77, 'Anish Kumar', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(100, 78, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Rejected'),
(101, 79, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(102, 80, 'K Mythri Vaishnavi', '953623104066', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(103, 81, 'Anish Kumar', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(104, 82, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(105, 83, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(106, 84, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(107, 85, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(108, 86, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(109, 87, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(110, 88, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(111, 89, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(112, 90, 'test', '953623104066', '3rd Year', 'CSE', 'B', 'Devi', 'Pending'),
(113, 91, 'Anish Kumar R', '953623104066', '3rd Year', 'CSE', 'B', 'Devi', 'Pending'),
(114, 92, 'Anish Kumar R', '953623104066', '3rd Year', 'CSE', 'B', 'Devi', 'Pending'),
(115, 93, 'Anish Kumar R', '953623104066', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(116, 94, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(117, 95, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(118, 96, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(119, 97, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(120, 98, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(121, 99, 'Mythri Vaishnavi', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Rejected'),
(122, 100, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(123, 101, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(124, 102, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(125, 103, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(126, 104, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(127, 105, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(128, 106, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Rejected'),
(129, 107, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Rejected'),
(130, 108, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Rejected'),
(131, 109, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Rejected'),
(132, 110, 'Saravanankamalesh', '953623104042', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Rejected'),
(133, 111, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(134, 112, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Rejected'),
(135, 113, 'Mythri Vaishnavi', '953623104066', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Rejected'),
(136, 114, 'Mythri Vaishnavi', '953623104066', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(137, 115, 'Mythri Vaishnavi', '953623104066', '3rd Year', 'CSE', 'B', 'Devi', 'Pending'),
(138, 116, 'Mythri Vaishnavi', '953623104066', '3rd Year', 'CSE', 'A', 'Vivek V', 'Pending'),
(139, 117, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(140, 118, 'Mythri Vaishnavi', '953623104066', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(141, 118, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Vivek V', 'Pending'),
(142, 119, 'Anish KuA', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Rejected'),
(143, 119, 'Mythri Vaishnavi', '953623104066', '3rd Year', 'CSE', 'B', 'Devi', 'Pending'),
(144, 120, 'Anish Kumar ', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(145, 120, 'Mythri Vaishnavi', '953623104066', '3rd Year', 'CSE', 'A', 'Vivek V', 'Accepted'),
(146, 120, 'David W Joshua', '953623104021', '3rd Year', 'CSE', 'B', 'Devi', 'Rejected'),
(147, 121, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(148, 122, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Vivek V', 'Rejected'),
(149, 123, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Rejected'),
(150, 124, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(151, 125, 'Diliban S', '953623104013', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(152, 126, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(153, 127, 'Aish Kumar', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(154, 128, 'Anish kUmar ', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(155, 129, 'Anish Kumar  R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(156, 130, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(157, 131, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(158, 132, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(159, 133, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(160, 134, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(161, 135, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(162, 136, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(163, 137, 'Kathir Vel', '953623104049', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(164, 138, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(165, 139, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(166, 140, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(167, 141, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(168, 142, 'Anish Kumar ', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(169, 143, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(170, 144, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(171, 144, 'Jeyaseelan', '953623104043', '3rd Year', 'CSE', 'A', 'Vivek V', 'Rejected'),
(172, 145, 'Anish kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(173, 146, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(174, 147, 'AK', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(175, 148, 'aNISH kUMAR ', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(176, 149, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(177, 150, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(178, 151, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(179, 151, 'Jeyaseelan', '953623104043', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(180, 152, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(181, 153, 'Anish Kumatr R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(182, 154, 'Anish Kumar', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(183, 155, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(184, 156, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(185, 157, 'Anish Kumar R', '953623104007', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(186, 158, 'Anish Kumar R', '953623104007', '3', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(187, 159, 'Anish Kumar R', '953623104007', '3', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(188, 160, 'Anish Kumar R', '953623104007', '3', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(189, 161, 'Anish Kumar R', '953623104007', '3', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(190, 162, 'Anish Kumar R', '953623104007', '3', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(191, 163, 'Anish Kumar R', '953623104007', '3', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(192, 164, 'Anish Kumar R', '953623104007', '3', 'CSE', 'A', 'Swarna Sudha', 'Pending'),
(193, 165, 'Anish Kumar R', '953623104007', '3', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(194, 166, 'David Joshua', '953623104021', '3', 'CSE', 'B', 'Vijay Amala Devi', 'Accepted'),
(195, 166, 'Lalith Krishnan', '953623104053', '3rd Year', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(196, 167, 'Anish Kumar R', '953623104007', '3', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(197, 167, 'David W Joshua', '953623104021', '3rd Year', 'CSE', 'B', 'Vijay Amala Devi', 'Pending'),
(198, 168, 'Anish Kumar R', '953623104007', '3', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(199, 168, 'David W Joshua', '953623104021', '3rd Year', 'CSE', 'B', 'Vijay Amala Devi', 'Accepted'),
(200, 169, 'Anish Kumar R', '953623104007', '3', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(201, 169, 'Joshua', '953623104020', '3rd Year', 'CSE', 'B', 'Vijay Amala Devi', 'Rejected'),
(202, 170, 'Lalith Krishna V M', '953623104053', '3', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(203, 170, 'Anish Kumar', '953623104007', '3rd Year', 'CSE', 'B', 'Vijay Amala Devi', 'Accepted'),
(204, 171, 'Anish Kumar R', '953623104007', '3', 'CSE', 'A', 'Swarna Sudha', 'Accepted'),
(205, 172, 'David W Joshua', '953623104021', '3', 'CSE', 'B', 'Vijay Amala Devi', 'Accepted'),
(206, 173, 'Afridha Fathima Syed Ibrahim', '953623104003', '3', 'CSE', 'B', 'Vijay Amala Devi', 'Accepted'),
(207, 174, 'David W Joshua', '953623104021', '3', 'CSE', 'B', 'Vijay Amala Devi', 'Pending'),
(208, 175, 'Janaramji R', '953623104040', '3', 'CSE', 'A', 'Vivek V', 'Pending'),
(209, 176, 'David W Joshua', '953623104021', '3', 'CSE', 'B', 'Vijay Amala Devi', 'Pending'),
(210, 177, 'David W Joshua', '953623104021', '3', 'CSE', 'B', 'Vijay Amala Devi', 'Pending'),
(211, 178, 'Anish Kumar R', '953623104007', '3', 'CSE', 'A', 'Swarna Sudha', 'Accepted');

-- --------------------------------------------------------

--
-- Table structure for table `principals`
--

CREATE TABLE `principals` (
  `id` int(11) NOT NULL,
  `register_no` varchar(50) NOT NULL,
  `name` varchar(100) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `principals`
--

INSERT INTO `principals` (`id`, `register_no`, `name`, `password`, `email`) VALUES
(1, '1032', 'Ganesan', 'c55b3f6a81a515bffa332d727560f1b5', 'anish.kmr.eng@gmail.com');

-- --------------------------------------------------------

--
-- Table structure for table `students`
--

CREATE TABLE `students` (
  `id` int(11) NOT NULL,
  `register_no` varchar(50) NOT NULL,
  `name` varchar(100) NOT NULL,
  `password` varchar(255) NOT NULL,
  `department` varchar(50) NOT NULL,
  `year` int(11) DEFAULT NULL,
  `section` varchar(10) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `students`
--

INSERT INTO `students` (`id`, `register_no`, `name`, `password`, `department`, `year`, `section`, `email`) VALUES
(1, '953623104125', 'Vignesh A P', 'fa78dccfca401d6eaf345c7d10969ea3', 'CSE', 3, 'A', ''),
(2, '953623104040', 'Janaramji R', 'b01e4afd11aeb2fef867d50dfbabf363', 'CSE', 3, 'A', ''),
(3, '953623104086', 'Ram Kumar S', '6b4efbb7d4671f098b2b7a89dec7392a', 'CSE', 3, 'A', ''),
(4, '953623104021', 'David W Joshua', 'e754cf95092440d08e7cff2465b6d0fa', 'CSE', 3, 'B', ''),
(5, '953623104100', 'Senthur Pramma', '8e86cb77cd25e8eba58152e285707b46', 'CSE', 3, 'B', ''),
(6, '953623104108', 'Siva Sakthi ', 'b8407b034c3dc0fbb5faf8316678750e', 'CSE', 3, 'B', ''),
(7, '953623104301', 'Darshan Sivakumar', '8c19b13d86082b3e5e75a9a6916950c9', 'CSE', 3, 'B', ''),
(8, '953623104070', 'Nithish Kumar B', '8a4a010098248c4494ca04005ddab346', 'CSE', 3, 'B', ''),
(9, '953623104032', 'Fadelullah M', '271d548ed2003f9cdd0c341014e41d1d', 'CSE', 3, 'B', ''),
(10, '953623104088', 'Rishikesh R', '02aff7d60126fd4011756fb2acdd5eaa', 'CSE', 3, 'B', ''),
(11, '953623104113', 'Shri Jaya Ganesh', 'a0739ece754cd3422336931128c98685', 'CSE', 3, 'B', ''),
(12, '953623104066', 'Mythri Vaishnavi', 'bf78840f13ef113f9cbd63f0e6aa4eab', 'CSE', 3, 'B', ''),
(13, '953623104003', 'Afridha Fathima Syed Ibrahim', '9ba443c3658ffa24e90a61c507091f59', 'CSE', 3, 'B', ''),
(14, '953623104062', 'Murugeswari R', '807f7d6a07dbc587a84fb3e921fe11f7', 'CSE', 3, 'B', ''),
(15, '953623104111', 'Sreevarthini M', '121db679f0c45cbb11c1c0d783b7fb9d', 'CSE', 3, 'B', ''),
(16, '953623104015', 'Asmitha V', '1a58445726f5d75ff6845f5dde34c087', 'CSE', 3, 'B', ''),
(17, '953623104014', 'Asmitha P K', '01e0f9532d6832db210b26330aac37f2', 'CSE', 3, 'B', ''),
(18, '953623104007', 'Anish Kumar R', 'c55b3f6a81a515bffa332d727560f1b5', 'CSE', 3, 'A', '');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `cas`
--
ALTER TABLE `cas`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `register_no` (`register_no`);

--
-- Indexes for table `hods`
--
ALTER TABLE `hods`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `register_no` (`register_no`);

--
-- Indexes for table `jas`
--
ALTER TABLE `jas`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `register_no` (`register_no`);

--
-- Indexes for table `lab_form`
--
ALTER TABLE `lab_form`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `lab_technicians`
--
ALTER TABLE `lab_technicians`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `register_no` (`register_no`);

--
-- Indexes for table `mentors`
--
ALTER TABLE `mentors`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `register_no` (`register_no`);

--
-- Indexes for table `od_applications`
--
ALTER TABLE `od_applications`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `od_team_members`
--
ALTER TABLE `od_team_members`
  ADD PRIMARY KEY (`id`),
  ADD KEY `od_id` (`od_id`);

--
-- Indexes for table `principals`
--
ALTER TABLE `principals`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `register_no` (`register_no`);

--
-- Indexes for table `students`
--
ALTER TABLE `students`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `register_no` (`register_no`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `cas`
--
ALTER TABLE `cas`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `hods`
--
ALTER TABLE `hods`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `jas`
--
ALTER TABLE `jas`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT for table `lab_form`
--
ALTER TABLE `lab_form`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `lab_technicians`
--
ALTER TABLE `lab_technicians`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `mentors`
--
ALTER TABLE `mentors`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=23;

--
-- AUTO_INCREMENT for table `od_applications`
--
ALTER TABLE `od_applications`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=179;

--
-- AUTO_INCREMENT for table `od_team_members`
--
ALTER TABLE `od_team_members`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=212;

--
-- AUTO_INCREMENT for table `principals`
--
ALTER TABLE `principals`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `students`
--
ALTER TABLE `students`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=19;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `od_team_members`
--
ALTER TABLE `od_team_members`
  ADD CONSTRAINT `od_team_members_ibfk_1` FOREIGN KEY (`od_id`) REFERENCES `od_applications` (`id`) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
